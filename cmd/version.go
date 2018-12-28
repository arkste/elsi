package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version represents the version
type Version struct {
	Version   string
	Commit    string
	BuildTime string
}

var version *Version

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number and exit (also --version)",
	Run:   printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// SetVersion sets the version
func SetVersion(v *Version) {
	version = v
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("elsi %s\n", version.Version)
}
