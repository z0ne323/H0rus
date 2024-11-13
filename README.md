# H0rus

*H0rus knows everything, H0rus sees everything*

![eye-of-horus-6078479_960_720](https://github.com/user-attachments/assets/308bd16d-64d3-4029-8385-4a2030f1babb)

*H0rus is a tool for cybersecurity researchers, penetration testers, or anyone who needs to scan a range of IPs for suspicious URLs and IOCs.*

## Overview

**H0rus** is a powerful network scanning tool designed to crawl IP ranges and search for Indicators of Compromise (IOCs) in HTTP and HTTPS URLs. 

Here are a few scenarios where this tool may be useful:
- **Incident response:** To find a webshell hosted across multiple web servers.
- **Threat hunting:** If you're looking for a specific IOC related to a C2 exposed on the internet.
- **And other applications...**

It uses ProtonVPN OpenVPN configuration files, ensuring traffic is securely routed through VPN tunnels during the scan for enhanced privacy.

## Features
- **IP Range Scanning**: Scans IP addresses between a start and end IP, excluding private ranges.
- **IOC Scanning**: Checks HTTP and HTTPS ports for the presence of specific IOCs (e.g., `malware.php`).
- **VPN Integration**: Uses ProtonVPN configuration files to route traffic securely.
- **Retry Mechanism**: Automatically retries the connection using different OpenVPN configurations if the initial connection fails.
- **Concurrent Scanning**: Uses goroutines to check multiple ports in parallel for efficient scanning (ports `80`, `443`, `8080` are set by default in `scanner.go`, but you can add more).
- **File Output**: Results are logged into a text file (`ioc_results.txt` by default, can be modified in `output.go`), located in the current working directory for later analysis.

## Setup

### Prerequisites

Before starting, you need to have the following installed:

- **Linux Distribution of your choice**

- **Go**
  - Debian/Ubuntu-based: `sudo apt install golang-go`
  - Fedora: `sudo dnf install go`
  - Arch: `sudo pacman -S go`

- **Git**
  - Debian/Ubuntu-based: `sudo apt install git`
  - Fedora: `sudo dnf install git`
  - Arch: `sudo pacman -S git`

- **OpenVPN**
  - Debian/Ubuntu-based: `sudo apt install openvpn`
  - Fedora: `sudo dnf install openvpn`
  - Arch: `sudo pacman -S openvpn`
  
- **ProtonVPN** [Sign up and set up your ProtonVPN account](https://protonvpn.com/)
  - After signing up, download ProtonVPN OpenVPN configuration files:
    1. Go [here](https://account.protonvpn.com/downloads#openvpn-configuration-files)
    2. Select: **GNU/Linux → TCP → Secure Core configs**
    3. Download all the configuration files to your `openvpn_confs` folder.

- **ProtonVPN Authentication File (`protonvpn_creds.txt`)**
  - Place your ProtonVPN credentials in this file. You can find the credentials [here](https://account.protonvpn.com/account-password#openvpn).
  - Set the correct permissions: `chmod 600 protonvpn_creds.txt`

- **Openresolv** (to avoid DNS leaks):
  - Debian/Ubuntu-based: `sudo apt install openresolv`
  - Fedora: `sudo dnf install openresolv`
  - Arch: `sudo pacman -S openresolv`

- **DNS Update Script from ProtonVPN**:
  1. Install `wget` if not already installed: `sudo apt install wget`
  2. Download the script:
     ```bash
     sudo wget "https://raw.githubusercontent.com/ProtonVPN/scripts/master/update-resolv-conf.sh" -O "/etc/openvpn/update-resolv-conf"
     ```
  3. Make the script executable:
     ```bash
     sudo chmod +x "/etc/openvpn/update-resolv-conf"
     ```


### Clone the Repository

Start by cloning the repository to your local machine:

```bash
git clone https://github.com/z0ne323/H0rus
cd H0rus
```

### Build the Tool
To build the tool from source, run:
```bash
go build -o H0rus ./cmd/main.go
```

### Usage
To run the tool, you need to provide the following values:

- **OpenVPN Authentication File** (`-auth`): Path to the ProtonVPN credentials file (e.g., protonvpn_creds.txt).
- **OpenVPN Configuration Directory** (`-config`): Path to the directory containing the .ovpn configuration files.
- **Indicator of Compromise (IOC)** (`-ioc`): The IOC you want to scan for, such as a URL or file name (e.g., malware.php).
- **Start IP and End IP Range** (`-startip` and `-endip`): Define the IP range you want to scan.

### You can provide these values in three different ways:

#### 1. Using Command-Line Flags

This is the most direct way to provide values when running the program. Here's the structure for using flags:

```bash
go run ./cmd/main.go -auth protonvpn_creds.txt -config openvpn_confs -ioc malware.php -startip 1.1.1.1 -endip 1.1.1.254
```

or (after building, see `Build the Tool` section for help) 

```bash
./H0rus -auth protonvpn_creds.txt -config openvpn_confs -ioc malware.php -startip 1.1.1.1 -endip 1.1.1.254
```

Explanation:
- `-auth protonvpn_creds.txt`: Specifies the path to your ProtonVPN authentication file.
- `-config openvpn_confs`: Specifies the directory containing your OpenVPN configuration files.
- `-ioc malware.php`: Specifies the IOC (Indicator of Compromise) to search for.
- `-startip 1.1.1.1`: The start IP address of the IP range to scan.
- `-endip 1.1.1.254`: The end IP address of the IP range to scan.

#### Using Interactive Input Prompts
Alternatively, you can run the tool without any flags. In this case, the program will ask for the necessary information interactively:

```bash
go run ./cmd/main.go
```

or (after building, with output example) 

```bash
# ./H0rus                                                                                                        
Enter the path to your authentication file for OpenVPN: protonvpn_creds.txt
Enter the path to your OpenVPN configuration file(s) directory: openvpn_confs
Enter your IOC: tmpukmfp.php
Enter the start IP address: 1.1.1.1
Enter the end IP address: 1.1.1.5
[+] IP range is valid!
[*] Starting our crawler...
[*] Testing IP: 1.1.1.1, with real_openvpn_confs/ireland-switzerland-01.protonvpn.tcp.ovpn
[*] Testing IP: 1.1.1.2, with real_openvpn_confs/netherlands-switzerland-01.protonvpn.tcp.ovpn
[*] Testing IP: 1.1.1.3, with real_openvpn_confs/unitedkingdom-sweden-01.protonvpn.tcp.ovpn
[*] Testing IP: 1.1.1.4, with real_openvpn_confs/france-switzerland-01.protonvpn.tcp.ovpn
[*] Testing IP: 1.1.1.5, with real_openvpn_confs/unitedstates-switzerland-02.protonvpn.tcp.ovpn
[*] No IOCs found during the scan.
[*] All IPs have been crawled.
```

Once you run the program, it will prompt you to enter the required values:
- **Enter the path to your authentication file for OpenVPN:** The path to the `protonvpn_creds.txt` file.
- **Enter the path to your OpenVPN configuration file(s) directory:** The path to the directory where the `.ovpn` files are stored.
- **Enter your IOC:** The IOC you're after, such as `malware.php`.
- **Enter the start IP address:** The starting IP address of the range to scan (e.g., `1.1.1.1`).
- **Enter the end IP address:** The ending IP address of the range to scan (e.g., `1.1.1.254`).

After entering all the values, the tool will proceed with the scan based on the provided inputs.

#### 3. Using a Configuration File (This is the default method if no flags or interactive input are provided.)
If you'd prefer to automate the configuration process, you can use a `JSON` configuration file. This eliminates the need to specify flags or enter values interactively each time.

Example configuration file (`config.json`):
```json
{
  "authFile": "/path/to/protonvpn_creds.txt",
  "configDir": "/path/to/openvpn_confs/",
  "ioc": "malware.php",
  "startIP": "1.1.1.1",
  "endIP": "1.1.1.254"
}
```

#### Important:
- The configuration file must be named `config.json` and stored in the root directory of this project. If it exists, the tool will use its values by default unless you override them using command-line flags or input prompts.
- If you provide any flags (e.g., `-auth`, `-config`, `-ioc`, etc.) or use the interactive mode, these will override any values from the `config.json` file.

### Missing Values & Prompts
While the configuration file and flags provide a convenient way to set the values, **it's not recommended** to leave out any required values. 

However, if any value is missing (whether from the command line, or the configuration file), the tool will automatically prompt you for the missing value.

For example:

- If you forget to provide the OpenVPN credentials file (`-auth`) via the command-line or in the `config.json` file, the tool will ask you to enter the path interactively.
- If you forget the `-ioc` flag or its value, the program will prompt you to specify which IOC you want to search for.

### Improvement / Ideas

- To increase the chances of successful connections, you could query an API to determine the geographical location of the target IP, allowing you to select a VPN server from a corresponding country (e.g., to access a Russian server, a Russian IP address is needed). This would be particularly useful for targets like `mil.ru`, which is only accessible from Russian IPs. 
- It might also be helpful to verify DNS leak protection when running the tool by querying an API such as `dnsleaktest.com`.
- In the future, adding support for other VPN service providers could further expand the tool's versatility.
- Find any ways possible to optimize the program / increase the speed 

### Contributing
Feel free to fork the repository and submit a pull request. 

It might be fun to see this project supporting other VPN service provider...
