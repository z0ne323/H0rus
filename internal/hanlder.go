package internal

import (
	"fmt"
	"math/rand"
)

// ConnectAndFetchIPResults connects to a remote IP using a VPN configuration and retrieves results for the specified IOC.
// The function will retry with a new OpenVPN configuration if the previous one fails.
func ConnectAndFetchIPResults(openVPNConfigFiles []string, openVPNAuthFile string, ip string, ioc string, maxRetries int) error {
	var err error
	for retries := 0; retries < maxRetries; retries++ {
		// Randomly select a config file from the shuffled list
		randomConfig := openVPNConfigFiles[rand.Intn(len(openVPNConfigFiles))]

		// Try to start OpenVPN with the selected configuration
		vpnProcess, err := StartOpenVPN(randomConfig, openVPNAuthFile)
		if err != nil {
			// If it fails, log the error and try again with a new config
			fmt.Printf("[-] Error starting OpenVPN with config %s: %v. Retrying...\n", randomConfig, err)
			continue
		}

		// If OpenVPN started successfully, defer killing the process and proceed
		defer vpnProcess.Process.Kill()

		// Let the user know which IP we test / what VPN conf file we currently use:
		fmt.Printf("[*] Testing IP: %s, with %s\n", ip, randomConfig)

		// Once connected, scan the HTTP ports for the specified IP and ioc
		scanHTTPPorts(ip, ioc)

		// If successful, break out of the retry loop
		return nil
	}

	// If all retries failed, return the last error encountered
	return fmt.Errorf("failed to establish VPN connection after %d retries: %v", maxRetries, err)
}
