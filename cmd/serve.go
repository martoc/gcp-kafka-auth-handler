package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API",
	Long:  `Start the API server`,
	Run: func(_ *cobra.Command, _ []string) {
		log.Default().Printf("Version: v%s", CLIVersion)
	},
}
