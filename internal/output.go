package internal

import (
	"fmt"
	"time"
)

// writeResultsToFile writes the IOC results to a file using the WriteToFile helper.
func WriteResultsToFile(results []string) {
	// Create the results file if it doesn't exist and write a timestamp
	fileName := "ioc_results.txt"
	header := fmt.Sprintf("[*] IOC Scan completed at %s\n", time.Now().Format(time.RFC3339))

	// Write the header and each IOC to the file
	if err := WriteToFile(fileName, header); err != nil {
		fmt.Println("[-] Error writing to file:", err)
		return
	}

	for _, result := range results {
		if err := WriteToFile(fileName, result); err != nil {
			fmt.Println("[-] Error writing IOC to file:", err)
		}
	}

	fmt.Printf("[+] IOCs found! Check the results in '%s'.", fileName)
}
