package cmd

import (
	"log"
	"os"

	"github.com/martoc/kafka-auth-handler/handler"
	"github.com/spf13/cobra"
)

const defaultPort = 14293

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API",
	Long: `Start the API server.

Environment variables:
  PROVIDER - Cloud provider for authentication: "gcp" (default) or "aws"
  REGION   - AWS region for MSK IAM authentication (required for AWS provider)`,
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("Version: v", CLIVersion)

		provider := os.Getenv("PROVIDER")
		region := os.Getenv("REGION")

		if provider == "" {
			provider = handler.ProviderGCP
		}

		log.Printf("Using provider: %s", provider)

		if provider == handler.ProviderAWS && region == "" {
			log.Fatal("REGION environment variable is required for AWS provider")
		}

		handler.StartServer(provider, region, defaultPort)
	},
}
