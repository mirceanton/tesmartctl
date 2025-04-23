package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sendRawCmd sends a raw hex command
var sendRawCmd = &cobra.Command{
	Use:   "raw [hex-command]",
	Short: "Send a raw hexadecimal command",
	Long: `Send raw hexadecimal commands directly to the KVM switch.
	This is useful for experimenting with the KVM protocol or for
	implementing commands that aren't yet supported by the CLI.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")
		address := net.JoinHostPort(ip, port)

		// Convert hex string to bytes
		hexCmd := args[0]
		payloadBytes, err := hex.DecodeString(hexCmd)
		if err != nil {
			return fmt.Errorf("invalid hex command: %v", err)
		}
		fmt.Printf("Sending %s to %s:%s...\n", hexCmd, ip, port)

		// Establish TCP connection
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			return fmt.Errorf("connection failed: %v", err)
		}
		defer conn.Close()
		if debug {
			fmt.Println("Connected successfully")
		}

		// Send the command
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Write(payloadBytes)
		if err != nil {
			return fmt.Errorf("error sending command: %v", err)
		}

		// Read response
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		// Buffer for response
		respBuf := make([]byte, 64)
		n, err := conn.Read(respBuf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading response: %v", err)
		}

		if n > 0 {
			// Convert response to hex
			respHex := hex.EncodeToString(respBuf[:n])
			fmt.Printf("Response: %s\n", respHex)

			// Also print as bytes if in debug mode
			if debug {
				fmt.Printf("Bytes: %v\n", bytes.NewBuffer(respBuf[:n]).Bytes())
			}
		} else {
			fmt.Println("No response received")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendRawCmd)
}
