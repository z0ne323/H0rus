package internal

import "flag"

// ParseFlags defines and parses command-line flags.
func ParseFlags() (string, string, string, string, string) {
	openVPNAuthFileFlag := flag.String("auth", "", "Path to your authentication file used with OpenVPN")
	openVPNConfigFilesDirectoryFlag := flag.String("config", "", "Path to your directory containing .ovpn configuration file(s)")
	iocFlag := flag.String("ioc", "", "IOC resource for URLs (e.g., tmpukmfp.php)")
	startIPFlag := flag.String("startip", "", "Start IP address for crawling (e.g., 1.1.1.1)")
	endIPFlag := flag.String("endip", "", "End IP address for crawling (e.g., 1.1.1.254)")
	flag.Parse()
	return *openVPNAuthFileFlag, *openVPNConfigFilesDirectoryFlag, *iocFlag, *startIPFlag, *endIPFlag
}

// GetOpenVPNAuthFile returns the path to the authentication file after validating or prompting.
func GetOpenVPNAuthFile(openVPNAuthFileFlag string, config *Config) (string, error) {
	if openVPNAuthFileFlag != "" {
		return openVPNAuthFileFlag, nil
	}
	if config != nil && config.OpenVPNAuthFile != "" {
		return config.OpenVPNAuthFile, nil
	}
	return PromptForOpenVPNAuthFile()
}

// GetOpenVPNConfigFilesDirectory returns the path to the config directory after validation or prompt.
func GetOpenVPNConfigFilesDirectory(openVPNConfigFilesDirectory string, config *Config) (string, error) {
	if openVPNConfigFilesDirectory != "" {
		return openVPNConfigFilesDirectory, nil
	}
	if config != nil && config.OpenVPNConfigFilesDirectory != "" {
		return config.OpenVPNConfigFilesDirectory, nil
	}
	return PromptForOpenVPNConfigFilesDirectory()
}

// GetIoc returns the IOC resource after validation or prompt.
func GetIoc(ioc string, config *Config) (string, error) {
	if ioc != "" {
		return ioc, nil
	}
	if config != nil && config.Ioc != "" {
		return config.Ioc, nil
	}
	return PromptForIoc()
}

// GetIP returns the IP address for crawling, either from the flag, config file, or prompt.
func GetIP(ip string, configIP string, ipType string) (string, error) {
	if ip != "" {
		return ip, nil
	}
	if configIP != "" {
		return configIP, nil
	}
	return PromptForIP(ipType)
}
