package cmd

import (
	"log"

	"github.com/martoc/gcp-kafka-auth-handler/handler"
	"github.com/spf13/cobra"
)

const defaultPort = 14293

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API",
	Long:  `Start the API server`,
	Run: func(cmd *cobra.Command, _ []string) {
		log.Default().Printf("Version: v%s", CLIVersion)
		handler.StartServer(defaultPort)
	},
}
