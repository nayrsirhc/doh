package doh

import (
	"github.com/spf13/cobra"
    "os"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
    Use:  "doh [command] [domain]",
    Short: "doh - a simple CLI to resolve DOH",
    Long: `DOH can be used to resolve DNS over HTTPS

This will be useful in the case you need to resolve external DNS for a domain but external resolution over port 53 is blocked`,
    Version: version,
    PreRunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            cmd.Help()
            os.Exit(0)
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
    },
}

func Execute() error {
    return rootCmd.Execute()
}