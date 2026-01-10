package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
}

var rootCmd = &cobra.Command{
	Use:   "kafka-auth-handler",
	Short: "Kafka OAuth2 token handler for GCP and AWS",
	Long:  `A lightweight HTTP server that provides OAuth2 tokens for Kafka clients authenticating against GCP or AWS.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
