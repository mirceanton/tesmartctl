package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool

	rootCmd = &cobra.Command{
		Use:   "tesmartctl",
		Short: "Control a TeSmart KVM switch",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/tesmartctl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if .config directory exists and create it if not
		configDir := filepath.Join(home, ".config")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "tesmartctl.yaml"
		viper.AddConfigPath(configDir)
		viper.SetConfigName("tesmartctl")
		viper.SetConfigType("yaml")
	}

	// Set default values
	viper.SetDefault("ip_address", "192.168.0.10")
	viper.SetDefault("port", "5000")

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Config file found
		if debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		// Config file not found, creating one with defaults
		if err := viper.SafeWriteConfig(); err != nil {
			fmt.Println("Failed to create config file:", err)
		} else if debug {
			fmt.Println("Created new config file:", viper.ConfigFileUsed())
		}
	}
}
