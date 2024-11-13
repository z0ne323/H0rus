package internal

import "fmt"

// PromptForOpenVPNAuthFile prompts the user for the authentication file path for OpenVPN.
func PromptForOpenVPNAuthFile() (string, error) {
	fmt.Print("Enter the path to your authentication file for OpenVPN: ")
	var openVPNAuthFile string
	if _, err := fmt.Scanln(&openVPNAuthFile); err != nil {
		return "", fmt.Errorf("[-] Error reading input: %w", err)
	}
	return openVPNAuthFile, nil
}

// PromptForOpenVPNConfigFilesDirectory prompts the user for the OpenVPN configuration files directory.
func PromptForOpenVPNConfigFilesDirectory() (string, error) {
	fmt.Print("Enter the path to your OpenVPN configuration file(s) directory: ")
	var openVPNConfigFilesDirectory string
	if _, err := fmt.Scanln(&openVPNConfigFilesDirectory); err != nil {
		return "", fmt.Errorf("[-] Error reading input: %w", err)
	}
	return openVPNConfigFilesDirectory, nil
}

// PromptForIoc prompts the user for the IOC Resource they want to check in the URLs.
func PromptForIoc() (string, error) {
	fmt.Print("Enter your IOC: ")
	var ioc string
	if _, err := fmt.Scanln(&ioc); err != nil {
		return "", fmt.Errorf("[-] Error reading input: %w", err)
	}
	return ioc, nil
}

// PromptForIP prompts the user for an IP address.
func PromptForIP(ipType string) (string, error) {
	fmt.Printf("Enter the %s IP address: ", ipType)
	var ip string
	if _, err := fmt.Scanln(&ip); err != nil {
		return "", fmt.Errorf("[-] Error reading input: %w", err)
	}
	return ip, nil
}
