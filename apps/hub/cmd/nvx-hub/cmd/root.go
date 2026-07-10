package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
	jsonOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "nvx-hub",
	Short: "Nevarix Hub service and CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", envOrDefault("NVX_CONFIG", "/etc/nvx/hub.yaml"), "Path to hub config file")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Emit JSON output")
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(serveCmd)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show hub running state and connection summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonOutput {
			fmt.Println(`{"hub_name":"hub-1","status":"stopped","managers_connected":0,"agents_connected":0}`)
			return nil
		}
		fmt.Println("Hub: hub-1")
		fmt.Println("Version: 0.1.0")
		fmt.Println("Status: stopped")
		fmt.Println("Managers connected: 0")
		fmt.Println("Agents connected: 0")
		fmt.Printf("Config: %s\n", configPath)
		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the hub gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("nvx-hub serve starting (config=%s)\n", configPath)
		fmt.Println("gRPC server placeholder — implement in Phase 2")
		select {}
	},
}
