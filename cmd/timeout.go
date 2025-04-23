package cmd

import (
	"fmt"
	"strings"

	"github.com/mirceanton/tesmartctl/internal/tesmart"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// timeoutCmd represents the timeout command
var timeoutCmd = &cobra.Command{
	Use:   "timeout [short|10|long|30|never|0|off]",
	Short: "Set the LED timeout on the KVM switch",
	Long: `Set the LED timeout on the TeSmart KVM switch.

Options:
  short, 10      : Set timeout to 10 seconds
  long, 30       : Set timeout to 30 seconds
  never, off, 0  : Disable LED timeout (LEDs always on)

Examples:
  tesmartctl timeout 10     # Set timeout to 10 seconds
  tesmartctl timeout never  # Disable timeout (LEDs always on)
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")
		arg := strings.ToLower(args[0])

		var command string
		var actionText string

		// Determine the command based on the argument
		switch arg {
		case "short", "10":
			command = "aabb03030aee"
			actionText = "set to 10 seconds"
		case "long", "30":
			command = "aabb03031eee"
			actionText = "set to 30 seconds"
		case "never", "off", "0":
			command = "aabb030300ee"
			actionText = "disabled (LEDs always on)"
		default:
			return fmt.Errorf("invalid timeout value: %s (use 'short'/'long'/'never' or '10'/'30'/'0')", arg)
		}

		// Send the command - timeout commands don't expect a response like buzzer commands
		_, err := tesmart.SendCommand(ip, port, command, false, debug)
		if err != nil {
			return fmt.Errorf("command failed: %v", err)
		}

		fmt.Printf("LED timeout %s\n", actionText)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(timeoutCmd)
}
