package cmd

import (
	"fmt"
	"strings"

	"github.com/mirceanton/tesmartctl/internal/tesmart"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buzzerCmd represents the buzzer command
var buzzerCmd = &cobra.Command{
	Use:   "buzzer [mute|unmute|0|1]",
	Short: "Control the KVM buzzer",
	Long: `Control the buzzer sound on the TeSmart KVM switch.

Examples:
  tesmartctl buzzer mute    # Mute the buzzer
  tesmartctl buzzer unmute  # Unmute the buzzer
  tesmartctl buzzer 0       # Mute the buzzer
  tesmartctl buzzer 1       # Unmute the buzzer
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")

		// Normalize the argument
		arg := strings.ToLower(args[0])

		var command string
		var actionText string

		// Determine the command based on the argument
		switch arg {
		case "mute", "0":
			command = "aabb03020000ee"
			actionText = "muted"
		case "unmute", "1":
			command = "aabb03020100ee"
			actionText = "unmuted"
		default:
			return fmt.Errorf("invalid argument: %s (use 'mute'/'unmute' or '0'/'1')", arg)
		}

		_, err := tesmart.SendCommand(ip, port, command, false, debug)
		if err != nil {
			return fmt.Errorf("command failed: %v", err)
		}
		fmt.Printf("buzzer %s\n", actionText)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buzzerCmd)
}
