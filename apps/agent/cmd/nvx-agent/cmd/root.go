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
	Use:   "nvx-agent",
	Short: "Nevarix Agent service and CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", envOrDefault("NVX_CONFIG", "/etc/nvx/agent.yaml"), "Path to agent config file")
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
	Short: "Show agent running state and connection summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonOutput {
			fmt.Println(`{"agent_name":"agent-1","status":"stopped","hubs_connected":0,"cache_pending":0}`)
			return nil
		}
		fmt.Println("Agent: agent-1")
		fmt.Println("Version: 0.1.0")
		fmt.Println("Status: stopped")
		fmt.Println("Hubs connected: 0")
		fmt.Println("Cache pending: 0")
		fmt.Printf("Config: %s\n", configPath)
		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the agent gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("nvx-agent serve starting (config=%s)\n", configPath)
		fmt.Println("gRPC client placeholder — implement in Phase 2")
		select {}
	},
}
