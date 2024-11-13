package main

import (
	"fmt"

	"github.com/z0ne323/H0rus/internal"
)

func main() {
	// Ensure the program is run with root privileges
	if !internal.IsRoot() {
		fmt.Println("[-] This program must be run as root.")
		return
	}

	// Parse command-line flags
	openVPNAuthFileFlag, openVPNConfigFilesDirectoryFlag, iocFlag, startIPFlag, endIPFlag := internal.ParseFlags()

	// Load the default configuration (config.json)
	config, err := internal.LoadConfig()
	if err != nil {
		fmt.Println("[-] Error loading config:", err)
		return
	}

	// Retrieve authentication file used with OpenVPN
	openVPNAuthFile, err := internal.GetOpenVPNAuthFile(openVPNAuthFileFlag, config)
	if err != nil {
		fmt.Println("[-] Error getting OpenVPN auth file:", err)
		return
	}

	// Load authentication file to verify it
	if _, err := internal.LoadOpenVPNAuthFile(openVPNAuthFile); err != nil {
		fmt.Println("[-] Error loading OpenVPN auth file:", err)
		return
	}

	// Retrieve directory containing all OpenVPN configuration files
	openVPNConfigFilesDirectory, err := internal.GetOpenVPNConfigFilesDirectory(openVPNConfigFilesDirectoryFlag, config)
	if err != nil {
		fmt.Println("[-] Error getting OpenVPN config files directory:", err)
		return
	}

	// Check if the directory is valid
	if err := internal.CheckOpenVPNConfigFilesDirectoryPath(openVPNConfigFilesDirectory); err != nil {
		fmt.Println("[-] Error checking OpenVPN config files directory:", err)
		return
	}

	// Retrieve IOC
	ioc, err := internal.GetIoc(iocFlag, config)
	if err != nil {
		fmt.Println("[-] Error getting IOC:", err)
		return
	}

	// Retrieve IP used for starting the crawler
	startIP, err := internal.GetIP(startIPFlag, config.StartIP, "start")
	if err != nil {
		fmt.Println("[-] Error with start IP:", err)
		return
	}

	// Verify IP format
	if validStartIP, _ := internal.IsValidIP(startIP); !validStartIP {
		fmt.Println("[-] Invalid start IP address:", startIP)
		return
	}

	// Retrieve IP used to stop the crawler
	endIP, err := internal.GetIP(endIPFlag, config.EndIP, "end")
	if err != nil {
		fmt.Println("[-] Error with end IP:", err)
		return
	}

	// Verify IP format
	if validEndIP, _ := internal.IsValidIP(endIP); !validEndIP {
		fmt.Println("[-] Invalid end IP address:", endIP)
		return
	}

	// Making sure the startIP isn't after the endIP
	validIPRange, err := internal.IsValidIPRange(startIP, endIP)
	if !validIPRange {
		fmt.Println("[-] Invalid IP range:", err)
		return
	} else {
		fmt.Println("[+] IP range is valid!")
	}

	// Initialize the random number generator
	internal.SeedRandom()

	// Get shuffled OpenVPN config files list from the directory
	openVPNConfigFiles, err := internal.GetShuffledConfigFiles(openVPNConfigFilesDirectory)
	if err != nil {
		fmt.Println("[-] Error fetching shuffled OpenVPN config files:", err)
		return
	}

	// Start crawling IPs and running OpenVPN connections
	fmt.Println("[*] Starting our crawler...")
	iocResults, err := internal.CrawlIPs(startIP, endIP, openVPNConfigFiles, openVPNAuthFile, ioc)
	if err != nil {
		fmt.Println("[-] Error during crawling:", err)
		return
	}

	// Handle results
	if len(iocResults) == 0 {
		// No IOCs found
		fmt.Println("[*] No IOCs found during the scan.")
	} else {
		// IOCs found, write them to a file and print a message
		internal.WriteResultsToFile(iocResults)
	}

	fmt.Println("[*] All IPs have been crawled.")
}
