package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Colors
const (
	NC     = "\033[0m"    // No Color
	RED    = "\033[0;31m" // Regular red
	GREEN  = "\033[0;32m" // Regular green
	YELLOW = "\033[1;33m" // Bold yellow
	BLUE   = "\033[0;34m" // Regular blue
	PURPLE = "\033[0;35m" // Regular purple
	CYAN   = "\033[0;36m" // Regular cyan
	WHITE  = "\033[37m"   // White
	BOLD   = "\033[1m"    // Bold text
	DIM    = "\033[2m"    // Dim text
)

var (
	outputFile *os.File
	verbose    bool
	threads    int
	mu         sync.Mutex
)



func showBanner() {
	banner := `
	██████╗  ██████╗ ██████╗  ██████╗ ██╗  ██╗██╗   ██╗███╗   ██╗████████╗
	██╔══██╗██╔═══██╗██╔══██╗██╔═══██╗██║  ██║██║   ██║████╗  ██║╚══██╔══╝
	██████╔╝██║   ██║██████╔╝██║   ██║███████║██║   ██║██╔██╗ ██║   ██║   
	██╔══██╗██║   ██║██╔══██╗██║   ██║██╔══██║██║   ██║██║╚██╗██║   ██║   
	██║  ██║╚██████╔╝██████╔╝╚██████╔╝██║  ██║╚██████╔╝██║ ╚████║   ██║   
	╚═╝  ╚═╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   
	`
	
	fmt.Printf("%s%s%s%s\n", BOLD, CYAN, banner, NC)
	fmt.Printf("%s%s%s%s\n", BOLD, GREEN, centerText("Robots.txt Directory Scanner", 80), NC)
	fmt.Printf("%s%s%s%s\n", DIM, YELLOW, centerText("v1.0.0", 80), NC)
	fmt.Printf("%s%s%s\n", PURPLE, strings.Repeat("═", 80), NC)
	fmt.Printf("%s%s%s%s\n", BOLD, PURPLE, centerText("Created by: Psikoz-coder", 80), NC)
	fmt.Printf("%s[%s●%s] Status: %sREADY TO HUNT%s\n", WHITE, GREEN, WHITE, GREEN, NC)
	fmt.Printf("%s%s%s\n\n", PURPLE, strings.Repeat("═", 80), NC)
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

func showUsage() {
	fmt.Printf("%sUsage: %s -l subdomain_list.txt [-o output.txt] [-v] [-t threads]%s\n", BLUE, os.Args[0], NC)
	fmt.Printf("%sParameters:%s\n", YELLOW, NC)
	fmt.Printf("  -l: Subdomain list file (required)\n")
	fmt.Printf("  -o: Output file (optional)\n")
	fmt.Printf("  -v: Verbose mode (show redirect information)\n")
	fmt.Printf("  -t: Number of threads (default: 10)\n")
	fmt.Printf("  -h: Help\n")
	fmt.Printf("\n%sExample: %s -l subdomains.txt -o results.txt -v -t 20%s\n", YELLOW, os.Args[0], NC)
}

func logOutput(message string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Print(message)
}

func logResult(message string) {
	if outputFile != nil {
		mu.Lock()
		defer mu.Unlock()
		outputFile.WriteString(message + "\n")
	}
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

func getRobotsContent(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

type HTTPResponse struct {
	StatusCode    int
	FinalURL      string
	RedirectChain []string
	Error         error
}

func getHTTPStatusWithRedirects(url string) HTTPResponse {
	var redirectChain []string
	
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectChain = append(redirectChain, req.URL.String())
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return HTTPResponse{
			StatusCode:    0,
			FinalURL:      url,
			RedirectChain: redirectChain,
			Error:         err,
		}
	}
	defer resp.Body.Close()
	
	return HTTPResponse{
		StatusCode:    resp.StatusCode,
		FinalURL:      resp.Request.URL.String(),
		RedirectChain: redirectChain,
		Error:         nil,
	}
}

func processPath(subdomain, path string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Clean path
	cleanPath := strings.ReplaceAll(path, "*", "")
	cleanPath = strings.ReplaceAll(cleanPath, "?", "")
	
	// Build full URL
	var testURL string
	if strings.HasPrefix(cleanPath, "/") {
		testURL = subdomain + cleanPath
	} else {
		testURL = subdomain + "/" + cleanPath
	}
	
	// Check HTTP status code
	httpResp := getHTTPStatusWithRedirects(testURL)
	
	// Color coding based on status code
	var color string
	var statusText string
	
	switch httpResp.StatusCode {
	case 200:
		color = GREEN
		statusText = "[200]"
	case 301, 302, 303, 307, 308:
		color = YELLOW
		statusText = fmt.Sprintf("[%d]", httpResp.StatusCode)
	case 403:
		color = BLUE
		statusText = "[403]"
	case 404:
		color = RED
		statusText = "[404]"
	case 500, 502, 503:
		color = RED
		statusText = fmt.Sprintf("[%d]", httpResp.StatusCode)
	case 0:
		color = RED
		statusText = "[TIMEOUT/ERROR]"
	default:
		color = NC
		statusText = fmt.Sprintf("[%d]", httpResp.StatusCode)
	}
	
	// Basic output
	output := fmt.Sprintf("  %s%s %s%s", color, statusText, testURL, NC)
	
	// Add redirect information
	if verbose {
		if len(httpResp.RedirectChain) > 0 {
			output += fmt.Sprintf("\n    %s↳ Redirects: %s%s", CYAN, strings.Join(httpResp.RedirectChain, " → "), NC)
		}
		if httpResp.FinalURL != testURL {
			output += fmt.Sprintf("\n    %s↳ Final URL: %s%s", CYAN, httpResp.FinalURL, NC)
		}
	}
	
	// Show redirect info even in non-verbose mode (short format)
	if !verbose && len(httpResp.RedirectChain) > 0 {
		output += fmt.Sprintf(" %s→ %s%s", CYAN, httpResp.FinalURL, NC)
	}
	
	output += "\n"
	logOutput(output)
	
	// Save to file
	logResult(fmt.Sprintf("%s %s", statusText, testURL))
	if len(httpResp.RedirectChain) > 0 {
		logResult(fmt.Sprintf("  ↳ Final: %s", httpResp.FinalURL))
	}
}

func processSubdomain(subdomain string) {
	// Add HTTP/HTTPS protocol
	if !strings.HasPrefix(subdomain, "http://") && !strings.HasPrefix(subdomain, "https://") {
		subdomain = "https://" + subdomain
	}
	
	logOutput(fmt.Sprintf("%s[*] Scanning: %s%s\n", YELLOW, subdomain, NC))
	
	// Get robots.txt
	robotsURL := subdomain + "/robots.txt"
	robotsContent, err := getRobotsContent(robotsURL)
	
	if err != nil {
		logOutput(fmt.Sprintf("%s[-] Robots.txt not found or inaccessible: %s%s\n", RED, robotsURL, NC))
		return
	}
	
	logOutput(fmt.Sprintf("%s[+] Robots.txt found: %s%s\n", GREEN, robotsURL, NC))
	
	// Find Disallow and Allow lines
	re := regexp.MustCompile(`^(Disallow|Allow):\s*(.+)$`)
	lines := strings.Split(robotsContent, "\n")
	var paths []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			path := strings.TrimSpace(matches[2])
			path = strings.ReplaceAll(path, "*", "")
			if path != "" && path != "/" {
				paths = append(paths, path)
			}
		}
	}
	
	// Remove duplicates and sort
	paths = removeDuplicates(paths)
	sort.Strings(paths)
	
	if len(paths) == 0 {
		logOutput(fmt.Sprintf("%s[!] No directories found in robots.txt%s\n", YELLOW, NC))
		return
	}
	
	logOutput(fmt.Sprintf("%s[+] Testing discovered directories... (%d paths)%s\n", BLUE, len(paths), NC))
	
	// Process with thread pool
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threads)
	
	for _, path := range paths {
		wg.Add(1)
		semaphore <- struct{}{}
		
		go func(p string) {
			defer func() { <-semaphore }()
			processPath(subdomain, p, &wg)
		}(path)
	}
	
	wg.Wait()
}

func main() {
	showBanner()
	
	var subdomainFile, outputFileName string
	var showHelp bool
	
	flag.StringVar(&subdomainFile, "l", "", "Subdomain list file (required)")
	flag.StringVar(&outputFileName, "o", "", "Output file (optional)")
	flag.BoolVar(&verbose, "v", false, "Verbose mode (show redirect information)")
	flag.IntVar(&threads, "t", 10, "Number of threads (default: 10)")
	flag.BoolVar(&showHelp, "h", false, "Help")
	flag.Parse()
	
	if showHelp {
		showUsage()
		return
	}
	
	if subdomainFile == "" {
		fmt.Printf("%sError: You must specify a subdomain list with -l parameter!%s\n", RED, NC)
		fmt.Printf("%sFor help: %s -h%s\n", YELLOW, os.Args[0], NC)
		os.Exit(1)
	}
	
	// Check file existence
	if _, err := os.Stat(subdomainFile); os.IsNotExist(err) {
		fmt.Printf("%sError: File %s not found!%s\n", RED, subdomainFile, NC)
		os.Exit(1)
	}
	
	// Open output file
	if outputFileName != "" {
		var err error
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			fmt.Printf("%sError: Cannot create output file: %s%s\n", RED, err, NC)
			os.Exit(1)
		}
		defer outputFile.Close()
		logOutput(fmt.Sprintf("[+] Output file: %s\n", outputFileName))
	}
	
	logOutput(fmt.Sprintf("[+] Reading subdomain list: %s\n", subdomainFile))
	logOutput(fmt.Sprintf("[+] Thread count: %d\n", threads))
	if verbose {
		logOutput(fmt.Sprintf("[+] Verbose mode: %sENABLED%s\n", GREEN, NC))
	}
	logOutput("[+] Scanning robots.txt files...\n\n")
	
	// Read subdomain file
	file, err := os.Open(subdomainFile)
	if err != nil {
		fmt.Printf("%sError: Cannot read file: %s%s\n", RED, err, NC)
		os.Exit(1)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := strings.TrimSpace(scanner.Text())
		if subdomain != "" {
			processSubdomain(subdomain)
			logOutput("\n")
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("%sError: File reading error: %s%s\n", RED, err, NC)
		os.Exit(1)
	}
	
	logOutput(fmt.Sprintf("%s[+] Scanning completed!%s\n", GREEN, NC))
	
	if outputFileName != "" {
		fmt.Printf("%s[+] Results saved to: %s%s\n", BLUE, outputFileName, NC)
	}
}