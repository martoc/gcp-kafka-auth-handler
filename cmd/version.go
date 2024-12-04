package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	CLIVersion string
	CLIOs      string
	CLIArch    string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions`,
	Run: func(_ *cobra.Command, _ []string) {
		jsonData := struct {
			Version string `json:"version"`
			Os      string `json:"os"`
			Arch    string `json:"arch"`
		}{
			Version: CLIVersion,
			Os:      CLIOs,
			Arch:    CLIArch,
		}
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error marshaling version data:", err)
		}
		fmt.Fprintln(os.Stdout, string(jsonBytes))
	},
}
