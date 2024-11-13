package internal

import (
	"fmt"
	"net"
)

// Define private IP ranges
var privateRanges = []string{
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"224.0.0.0/4",
	"240.0.0.0/4",
}

// IsValidIP checks if the provided IP address is valid and not within a private range.
func IsValidIP(ip string) (bool, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, fmt.Errorf("invalid IP format: %s", ip)
	}

	// Check if the IP is in any of the private ranges
	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return false, fmt.Errorf("error parsing private range CIDR: %v", err)
		}

		if network.Contains(parsedIP) {
			return false, fmt.Errorf("IP %s is in a private range: %s", ip, cidr)
		}
	}

	// If the IP is valid and not private
	return true, nil
}

// IsValidIPRange validates that the start IP is less than or equal to the end IP
func IsValidIPRange(startIP string, endIP string) (bool, error) {
	// Parse the start and end IPs
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	// Check if either IP is invalid
	if start == nil {
		return false, fmt.Errorf("invalid start IP format: %s", startIP)
	}
	if end == nil {
		return false, fmt.Errorf("invalid end IP format: %s", endIP)
	}

	// Convert the IPs to big.Int for lexicographical comparison
	startBigInt := ipToBigInt(start)
	endBigInt := ipToBigInt(end)

	// Check if the start IP is greater than the end IP
	if startBigInt.Cmp(endBigInt) > 0 {
		return false, fmt.Errorf("start IP %s cannot be greater than end IP %s", startIP, endIP)
	}

	// If the IPs are valid and the start IP is less than or equal to the end IP
	return true, nil
}
