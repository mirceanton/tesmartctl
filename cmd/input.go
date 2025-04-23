package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mirceanton/tesmartctl/internal/tesmart"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// inputCmd represents the input command
var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Get or set the active input on the KVM switch",
	Long: `Get the current active input or switch to a specific input port on the TeSmart KVM switch.

Use "tesmartctl input get" to see the current active input.
Use "tesmartctl input set <number>" to switch to a specific input port.
`,
}

// inputGetCmd represents the input get subcommand
var inputGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current active input",
	Long:  `Get the current active input port on the TeSmart KVM switch.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")

		if debug {
			fmt.Printf("Getting current active input from KVM at %s:%s...\n", ip, port)
		}

		// Command to get the current input: 0xAA 0xBB 0x03 0x10 0x00 0xEE
		command := "aabb031000ee"

		// Send the command - this one expects a response
		response, err := tesmart.SendCommand(ip, port, command, true, debug)
		if err != nil {
			return fmt.Errorf("failed to get current input: %v", err)
		}

		// Parse the response to get the active port
		if len(response) < 12 {
			return fmt.Errorf("invalid response length: %s", response)
		}

		// The active port is in the 5th byte (10th and 11th characters in the hex string)
		portHex := response[8:10]

		portNum, err := strconv.ParseUint(portHex, 16, 8)
		if err != nil {
			return fmt.Errorf("failed to parse port number from response: %v", err)
		}

		// TeSmart port numbers are 0-based in the protocol
		portNum = portNum + 1

		fmt.Println(portNum)

		return nil
	},
}

// inputSetCmd represents the input set subcommand
var inputSetCmd = &cobra.Command{
	Use:   "set [port-number]",
	Short: "Set the active input",
	Long: `Set the active input port on the TeSmart KVM switch.

The port number should be between 1 and 16, representing PC1 through PC16.

Example:
  tesmartctl input set 3  # Switch to PC3
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")

		// Parse the port number from the argument
		portNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid port number: %s", args[0])
		}

		// Validate the port number
		if portNum < 1 || portNum > 16 {
			return fmt.Errorf("port number must be between 1 and 16, got: %d", portNum)
		}

		// Format the command: 0xAA 0xBB 0x03 0x01 [port] 0xEE
		command := fmt.Sprintf("aabb0301%02xee", portNum)

		if debug {
			fmt.Printf("Switching to input %d on KVM at %s:%s...\n", portNum, ip, port)
		}

		// Send the command - this one expects a response
		response, err := tesmart.SendCommand(ip, port, command, true, debug)
		if err != nil {
			return fmt.Errorf("failed to switch input: %v", err)
		}

		// Verify the response indicates success
		if len(response) < 12 {
			return fmt.Errorf("invalid response length: %s", response)
		}

		// For input switching, the KVM typically responds with 0xAA 0xBB 0x03 0x11 [new-port] 0xEE
		// Let's check if the response contains the expected pattern
		respPrefix := strings.HasPrefix(response, "aabb0311")

		if !respPrefix {
			if debug {
				fmt.Printf("Unexpected response format: %s\n", response)
			}
			fmt.Println("Switch command sent, but received unexpected response")
			return nil
		}

		// Extract the new active port from the response
		respPortHex := response[8:10]
		respPortNum, err := strconv.ParseUint(respPortHex, 16, 8)
		if err != nil {
			fmt.Println("Successfully switched input, but couldn't parse response details")
			return nil
		}

		// TeSmart port numbers are 0-based in the protocol
		respPortNum++

		fmt.Printf("switched to input %d\n", respPortNum)

		return nil
	},
}

func init() {
	inputCmd.AddCommand(inputGetCmd)
	inputCmd.AddCommand(inputSetCmd)
	rootCmd.AddCommand(inputCmd)
}
