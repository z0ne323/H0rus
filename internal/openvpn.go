package internal

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// monitorOutput reads the output of OpenVPN command (stdout or stderr) to check for specific messages.
func monitorOutput(pipe *bufio.Reader, connectionEstablished chan bool, successMsg string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		// Print each line for better debugging
		// fmt.Println("OpenVPN Output:", line)

		// Check for success message or error
		if strings.Contains(line, successMsg) {
			connectionEstablished <- true
			return
		}
	}
	connectionEstablished <- false
}

// StartOpenVPN initializes an OpenVPN connection using the provided config and authentication file.
func StartOpenVPN(configFile string, authFile string) (*exec.Cmd, error) {
	cmd := exec.Command("openvpn", "--config", configFile, "--auth-user-pass", authFile, "--auth-nocache")

	// Get stdout and stderr pipes for the OpenVPN process
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the OpenVPN process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start OpenVPN: %w", err)
	}

	// Create a channel to wait for connection status
	connectionEstablished := make(chan bool)

	// Monitor OpenVPN's stdout and stderr for connection success or failure
	go monitorOutput(bufio.NewReader(stdout), connectionEstablished, "Initialization Sequence Completed")
	go monitorOutput(bufio.NewReader(stderr), connectionEstablished, "Options error")

	// Wait for connection status or timeout
	select {
	case success := <-connectionEstablished:
		if success {
			// Uncomment for debugging if needed
			// fmt.Println("[+] OpenVPN launched successfully.")
			return cmd, nil
		}
		cmd.Process.Kill()
		return nil, fmt.Errorf("failed to establish OpenVPN connection")
	case <-time.After(3 * time.Second): // Timeout after 3 seconds
		cmd.Process.Kill()
		return nil, fmt.Errorf("failed to establish OpenVPN connection within timeout")
	}
}
