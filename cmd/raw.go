package cmd

import (
	"fmt"

	"github.com/mirceanton/tesmartctl/internal/tesmart"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sendRawCmd sends a raw hex command
var sendRawCmd = &cobra.Command{
	Use:   "raw [hex-command]",
	Short: "Send a raw hexadecimal command",
	Long: `Send raw hexadecimal commands directly to the KVM switch.
This is useful for experimenting with the KVM protocol or for
implementing commands that aren't yet supported by the CLI.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")
		command := args[0]

		if debug {
			fmt.Printf("Sending command %s to %s:%s...\n", command, ip, port)
		}

		response, err := tesmart.SendCommand(ip, port, command, debug)
		if err != nil {
			return fmt.Errorf("command failed: %v", err)
		}

		if response != "" {
			fmt.Printf("%s\n", response)
		} else {
			fmt.Println("No response received")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendRawCmd)
}
