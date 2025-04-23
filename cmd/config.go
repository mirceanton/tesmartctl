package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `View or modify configuration settings for the KVM connection.

Examples:
	tesmartctl config get     # Show current config
	tesmartctl config set ip 10.0.0.4  # Set KVM IP in config
	tesmartctl config set port 9000    # Set KVM port in config
`,
}

// getConfigCmd represents the config get subcommand
var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "View current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current configuration:\n")
		fmt.Printf("  IP Address: %s\n", viper.GetString("ip_address"))
		fmt.Printf("  Port: %s\n", viper.GetString("port"))
		fmt.Printf("\nConfiguration file: %s\n", viper.ConfigFileUsed())
	},
}

// setConfigCmd represents the config set subcommand
var setConfigCmd = &cobra.Command{
	Use:   "set [ip|port] [value]",
	Short: "Modify configuration",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		param := args[0]
		value := args[1]

		switch param {
		case "ip":
			viper.Set("ip_address", value)
			fmt.Printf("IP address set to: %s\n", value)
		case "port":
			viper.Set("port", value)
			fmt.Printf("Port set to: %s\n", value)
		default:
			return fmt.Errorf("unknown parameter: %s", param)
		}

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("error writing config: %v", err)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(getConfigCmd)
	configCmd.AddCommand(setConfigCmd)

	rootCmd.AddCommand(configCmd)
}
