package internal

import (
	"fmt"
	"math/big"
	"net"
)

// nextIP returns the next IP in the sequence as a *big.Int
func nextIP(ip *big.Int) *big.Int {
	ip.Add(ip, big.NewInt(1)) // Increment the IP by 1
	return ip
}

// CrawlIPs starts the scanning process for a range of IPs and returns any found IOCs.
func CrawlIPs(startIP string, endIP string, openVPNConfigFiles []string, openVPNAuthFile string, IOC string) ([]string, error) {
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	// Convert the start and end IPs to big.Int for comparison
	startBigIntIP := ipToBigInt(start)
	endBigIntIP := ipToBigInt(end)

	var results []string

	// Iterate through IPs from start to end (inclusive)
	for ip := startBigIntIP; ip.Cmp(endBigIntIP) <= 0; ip = nextIP(ip) {
		ipStr := bigIntToIP(ip).String() // Convert back to net.IP for printing

		// Call IsValidIP from validation.go
		if valid, err := IsValidIP(ipStr); valid {
			// IP is valid and not in a private range, proceed
			// Provide the maxRetries argument here (3 retries in this case)
			err := ConnectAndFetchIPResults(openVPNConfigFiles, openVPNAuthFile, ipStr, IOC, 3)
			if err != nil {
				fmt.Println("Error:", err)
			}

			// Collect the results from the scan
			ipResults := scanHTTPPorts(ipStr, IOC)
			results = append(results, ipResults...)
		} else {
			// IP is invalid or in a private range
			fmt.Printf("[*] Skipping IP: %s (Reason: %s)\n", ipStr, err.Error())
		}
	}
	return results, nil
}
