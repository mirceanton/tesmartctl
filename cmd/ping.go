package cmd

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Test connectivity to the KVM",
	Long:  `Test TCP connectivity to the KVM switch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := viper.GetString("ip_address")
		port := viper.GetString("port")
		address := net.JoinHostPort(ip, port)
		log.Infof("Testing TCP connection to KVM at %s...\n", address)

		log.Debug("Establish TCP connection with timeout")
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			fmt.Println("\nTroubleshooting suggestions:")
			fmt.Printf("  1. Verify KVM is powered on and connected to the network\n")
			fmt.Printf("  2. Check IP address configuration (current: %s)\n", ip)
			fmt.Printf("  3. Try to ping the IP address: ping %s\n", ip)
			fmt.Printf("  4. Try direct TCP connectivity: nc -zv %s %s\n", ip, port)
			fmt.Printf("  5. Update config if needed: tesmartctl config set ip <new-ip>\n")
			return err
		}
		defer conn.Close()

		fmt.Printf("✅ Successfully connected to KVM at %s:%s\n", ip, port)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
