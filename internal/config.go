package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	OpenVPNAuthFile             string `json:"openVPNAuthFile"`
	OpenVPNConfigFilesDirectory string `json:"openVPNConfigFilesDirectory"`
	Ioc                         string `json:"ioc"`
	StartIP                     string `json:"startIP"`
	EndIP                       string `json:"endIP"`
}

var configFilePath = "config.json" // Default location for the config file

// LoadConfig loads H0rus configuration from a JSON file.
func LoadConfig() (*Config, error) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file %s not found", configFilePath)
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

type AuthFileOpenVPNFormat struct {
	Username string
	Password string
}

// parseOpenVPNAuthFile reads the username and password from the given file.
func parseOpenVPNAuthFile(file *os.File, config *AuthFileOpenVPNFormat) error {
	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	if len(lines) < 2 {
		return fmt.Errorf("file must contain at least two lines (username and password)")
	}

	config.Username = strings.TrimSpace(lines[0])
	config.Password = strings.TrimSpace(lines[1])
	return nil
}

// checkDefaultCredentials ensures that the credentials are not set to default values.
func checkDefaultCredentials(config *AuthFileOpenVPNFormat) error {
	if config.Username == "username" || config.Password == "password" {
		return fmt.Errorf("username and/or password are set to default values.")
	}
	return nil
}

// LoadOpenVPNConfig reads the OpenVPN authentication file and populates the provided OpenVPNConfig instance
// Also checks for default creds
func LoadOpenVPNConfig(filename string, config *AuthFileOpenVPNFormat) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := parseOpenVPNAuthFile(file, config); err != nil {
		return err
	}

	return checkDefaultCredentials(config)
}

var random *rand.Rand

// SeedRandom initializes the random number generator.
func SeedRandom() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GetConfigFiles retrieves all the .ovpn files from the specified directory.
func GetConfigFiles(directory string) ([]string, error) {
	var configFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".ovpn" {
			configFiles = append(configFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return configFiles, nil
}

// GetShuffledConfigFiles retrieves and shuffles OpenVPN configuration files.
func GetShuffledConfigFiles(directory string) ([]string, error) {
	files, err := GetConfigFiles(directory)
	if err != nil {
		return nil, err
	}

	random.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})
	return files, nil
}
