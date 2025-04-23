package tesmart

import (
	"encoding/hex"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// SendCommand sends a raw hex command string to the KVM switch and returns the response as a hex string
// The input should be a hex string (without 0x prefix or spaces)
// If expectResponse is false, it won't wait for a response and will return empty string
func SendCommand(ipAddress, port, hexCommand string, expectResponse bool, debug bool) (string, error) {
	log.Debugf("Deconding command bytes...\n")
	cmdBytes, err := hex.DecodeString(hexCommand)
	if err != nil {
		return "", err
	}
	log.Debugf("Decoded: %v\n", cmdBytes)

	log.Debugf("Establishing TCP connection...\n")
	address := net.JoinHostPort(ipAddress, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	log.Debugf("Connection established!\n")

	log.Debugf("Sending %d bytes\n", len(cmdBytes))
	_, err = conn.Write(cmdBytes)
	if err != nil {
		return "", err
	}
	log.Debugf("Command sent successfully")

	// If we don't expect a response, return early
	if !expectResponse {
		log.Debugf("Not waiting for response as requested")
		return "", nil
	}

	log.Debugf("Reading response (expecting 6 bytes for TeSmart protocol)")
	respBuf := make([]byte, 6)
	n, err := io.ReadFull(conn, respBuf)
	if err != nil {
		return "", err
	}
	log.Debugf("Received %d bytes: %v\n", n, respBuf[:n])

	// Convert response to hex string
	respHex := hex.EncodeToString(respBuf[:n])
	log.Debugf("Response hex: %s\n", respHex)

	return respHex, nil
}
