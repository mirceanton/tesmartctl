package tesmart

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"time"
)

// SendCommand sends a raw hex command string to the KVM switch and returns the response as a hex string
// The input should be a hex string (without 0x prefix or spaces)
func SendCommand(ipAddress, port, hexCommand string, debug bool) (string, error) {
	// Convert hex string to bytes
	if debug {
		fmt.Print("Deconding command bytes...\n")
	}
	cmdBytes, err := hex.DecodeString(hexCommand)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Printf("Decoded: %v\n", cmdBytes)
	}

	// Establish TCP connection
	address := net.JoinHostPort(ipAddress, port)
	if debug {
		fmt.Printf("Connecting to %s\n", address)
	}
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if debug {
		fmt.Println("Connection established")
	}

	// Send the command
	if debug {
		fmt.Printf("Sending %d bytes\n", len(cmdBytes))
	}
	_, err = conn.Write(cmdBytes)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Println("Command sent successfully")
	}

	// Read response (expecting 6 bytes for TeSmart protocol)
	respBuf := make([]byte, 6)
	if debug {
		fmt.Println("Reading response...")
	}

	n, err := io.ReadFull(conn, respBuf)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Printf("Received %d bytes: %v\n", n, respBuf[:n])
	}

	// Convert response to hex string
	respHex := hex.EncodeToString(respBuf[:n])
	if debug {
		fmt.Printf("Response hex: %s\n", respHex)
	}

	return respHex, nil
}
