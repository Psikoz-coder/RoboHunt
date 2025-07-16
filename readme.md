# RoboHunt

```
        ██████╗  ██████╗ ██████╗  ██████╗ ██╗  ██╗██╗   ██╗███╗   ██╗████████╗
        ██╔══██╗██╔═══██╗██╔══██╗██╔═══██╗██║  ██║██║   ██║████╗  ██║╚══██╔══╝
        ██████╔╝██║   ██║██████╔╝██║   ██║███████║██║   ██║██╔██╗ ██║   ██║   
        ██╔══██╗██║   ██║██╔══██╗██║   ██║██╔══██║██║   ██║██║╚██╗██║   ██║   
        ██║  ██║╚██████╔╝██████╔╝╚██████╔╝██║  ██║╚██████╔╝██║ ╚████║   ██║   
        ╚═╝  ╚═╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   
        
                          Robots.txt Directory Scanner
                                     v1.0.0
```

**Created by: [Psikoz-coder](https://github.com/Psikoz-coder)**

---

`RoboHunt` is a high-performance, concurrent `robots.txt` scanner designed for security professionals and bug bounty hunters. It efficiently scans a list of subdomains to find accessible directories and paths listed in their `robots.txt` files, helping you uncover potential vulnerabilities and expand your attack surface.

## 🌟 Features

- **Concurrent Scanning**: Utilizes goroutines to scan multiple subdomains at once for maximum speed.
- **Verbose Mode**: Get detailed information about redirects and scan status.
- **Customizable Threads**: Adjust the number of concurrent threads to match your needs.
- **Flexible Output**: Save the results to a file for further analysis.
- **Clear & Organized**: Color-coded output for easy readability.

## 📦 Installation

### Prerequisites
- Go 1.16 or higher
- A list of subdomains to scan

### Install with go install

```bash
go install -v github.com/Psikoz-coder/RoboHunt@latest
```

### Install from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Psikoz-coder/RoboHunt.git
   cd RoboHunt
   ```

2. **Build the executable:**
   ```bash
   go build -o robohunt robohunt.go
   ```

3. **(Optional) Move to your bin directory:**
   ```bash
   sudo mv robohunt /usr/local/bin/
   ```

## 🚀 Usage

### Basic Command

To start scanning, you need a text file containing a list of subdomains (one per line).

```bash
robohunt -l subdomains.txt
```

### Saving the Output

Save the scan results to a file for later review.

```bash
robohunt -l subdomains.txt -o results.txt
```

### Verbose Mode

For more detailed output, including redirect information, use the `-v` flag.

```bash
robohunt -l subdomains.txt -v
```

### Adjusting Threads

Control the performance by setting the number of concurrent threads.

```bash
robohunt -l subdomains.txt -t 20
```

## 🛠️ Command-Line Parameters

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-l` | **(Required)** Path to the subdomain list file. | | 
| `-o` | Path to the output file. | `(none)` |
| `-v` | Enable verbose mode to show redirect info. | `false` |
| `-t` | Number of concurrent threads. | `10` |
| `-h` | Show the help message. | | 

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/Psikoz-coder/RoboHunt/issues).

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## ⚠️ Legal Disclaimer

This tool is intended for ethical and legal security testing purposes only. Usage of `RoboHunt` for any unauthorized or malicious activities is strictly prohibited. The author is not responsible for any misuse or damage caused by this program.

---

⭐ **If you find this tool useful, please give it a star!**