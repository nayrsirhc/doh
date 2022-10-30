package doh

import (
	"github.com/spf13/cobra"
	// "github.com/nayrsirhc/doh/pkg/doh"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
    Use:  "doh [command] [domain]",
    Short: "doh - a simple CLI to transform and inspect strings",
    Long: `doh is a super fancy CLI (kidding)

One can use doh to modify or inspect strings straight from the terminal`,
    Version: version,
    Run: func(cmd *cobra.Command, args []string) {
    },
}

func Execute() error {
    return rootCmd.Execute()
}