package cmd

import (
	"fmt"

	"github.com/mirceanton/tesmartctl/internal/tesmart"
	log "github.com/sirupsen/logrus"
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

Examples:
	tesmartctl raw aabb031000ee # read current active input
	tesmartctl raw aabb030200ee # mute buzzer

Reference: https://support.tesmart.com/hc/en-us/article_attachments/27712362494361
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")
		log.Debugf("Sending raw command to KVM at %s:%s...\n", ip, port)

		command := args[0]
		log.Debugf("Command to send %s\n", command)

		log.Infof("Sending command %s to %s:%s...", command, ip, port)
		response, err := tesmart.SendCommand(ip, port, command, true, debug)
		if err != nil {
			return fmt.Errorf("command failed: %v", err)
		}
		log.Debugf("Got response: %s", response)

		log.Debugf("Verifying the response...")
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
