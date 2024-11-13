package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// checkOpenVPNAuthFile checks if the file can be read and does not contain default credentials.
func checkOpenVPNAuthFile(filename string) error {
	var config AuthFileOpenVPNFormat
	if err := LoadOpenVPNConfig(filename, &config); err != nil {
		return err
	}

	return checkDefaultCredentials(&config)
}

// LoadOpenVPNAuthFile checks for the authentication file and validates it.
func LoadOpenVPNAuthFile(OpenVPNAuthFileFlag string) (string, error) {
	if !FileExists(OpenVPNAuthFileFlag) {
		return "", fmt.Errorf("auth file does not exist: %s", OpenVPNAuthFileFlag)
	}

	if err := checkOpenVPNAuthFile(OpenVPNAuthFileFlag); err != nil {
		return "", err
	}

	return OpenVPNAuthFileFlag, nil
}

// checkOpenVPNConfigFilesDirectoryContents checks for .ovpn files in the specified directory
func checkOpenVPNConfigFilesDirectoryContents(directory string) error {
	hasOVPN := false
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".ovpn" {
			hasOVPN = true
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	if !hasOVPN {
		return fmt.Errorf("no .ovpn files found in the specified directory: %s", directory)
	}

	return nil
}

// CheckOpenVPNConfigFilesDirectoryPath checks if the configuration directory is valid and contains .ovpn files.
func CheckOpenVPNConfigFilesDirectoryPath(directory string) error {
	if directory == "" {
		return fmt.Errorf("no configuration directory specified.")
	}
	if err := checkDirectory(directory); err != nil {
		return err
	}
	if err := checkOpenVPNConfigFilesDirectoryContents(directory); err != nil {
		return err
	}
	return nil
}
