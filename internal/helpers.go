package internal

import (
	"fmt"
	"math/big"
	"net"
	"os"
	"sync"
)

// IsRoot checks if the program is running with root privileges.
func IsRoot() bool {
	return os.Geteuid() == 0
}

// FileExists checks if a file exists at the given filename path.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// checkDirectory validates if the directory exists and is a directory.
func checkDirectory(directory string) error {
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return fmt.Errorf("the specified path does not exist: %s", directory)
	} else if !info.IsDir() {
		return fmt.Errorf("the specified path is not a directory: %s", directory)
	}
	return nil
}

// Convert an IP address to a big.Int
func ipToBigInt(ip net.IP) *big.Int {
	// Pad the IP to make sure it's a 16-byte slice for both IPv4 and IPv6
	ip = ip.To16()
	return new(big.Int).SetBytes(ip)
}

// Convert a big.Int back to a net.IP
func bigIntToIP(bigInt *big.Int) net.IP {
	// Convert back to 16-byte array and return as net.IP
	ipBytes := bigInt.Bytes()
	ip := make(net.IP, 16)
	copy(ip[16-len(ipBytes):], ipBytes)
	return ip
}

// Declare a mutex to protect access to file writing
var fileMutex sync.Mutex

// WriteToFile safely writes data to the specified file.
// It uses a mutex to ensure that only one goroutine writes to the file at a time.
func WriteToFile(filename string, data string) error {
	// Lock the mutex to ensure that only one goroutine writes at a time
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Open the file in append mode
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.WriteString(data + "\n") // Add newline to separate entries
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	return nil
}
