package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Define a list of ports to check
var ports = []int{80, 443, 8080}

// Declare a mutex and a slice to store IOC results
var iocResults = struct {
	sync.Mutex
	Results []string
}{}

// checkHTTPPort attempt to reach the exact IP, port, ioc provided and also handle redirects
func checkHTTPPort(ip string, port int, ioc string, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when the goroutine is done

	// Construct the URL with the provided IP, port, and IOC signature
	var url string
	if port == 443 {
		url = fmt.Sprintf("https://%s/%s", ip, ioc) // HTTPS for port 443
	} else {
		url = fmt.Sprintf("http://%s:%d/%s", ip, port, ioc) // HTTP for other ports
	}

	// Set a timeout for the request
	client := &http.Client{
		Timeout: 1 * time.Second, // Timeout after 1 seconds
		// Disable redirect following
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// If a redirect happens, simply don't follow it
			return http.ErrUseLastResponse
		},
	}

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {
		// If there's an error (e.g., the port is closed or unreachable), just return without processing
		return
	}
	defer resp.Body.Close()

	// Only process if we get a 200 OK response
	if resp.StatusCode == 200 {
		// Lock the iocResults mutex and append the IOC to the results slice
		iocResults.Lock()
		iocResults.Results = append(iocResults.Results, fmt.Sprintf("[+] IOC found!! => %s:%d/%s", ip, port, ioc))
		iocResults.Unlock()
	} else {
		// Ignore all other status codes (such as 301, 404, etc.) and do not log them
		return
	}
}

// scanHTTPPorts now checks only the exact IP:Port and IOC you want
func scanHTTPPorts(ip string, ioc string) []string {
	// Use sync.WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Launch a goroutine for each port to check it concurrently
	for _, port := range ports {
		wg.Add(1) // Increment the WaitGroup counter
		go checkHTTPPort(ip, port, ioc, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// After crawling, lock the iocResults mutex and return results
	iocResults.Lock()
	defer iocResults.Unlock()

	// Return a copy of the results to main.go
	return append([]string{}, iocResults.Results...) // Returning a copy to avoid external mutation
}
