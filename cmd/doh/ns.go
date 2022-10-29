package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var nsCmd = &cobra.Command{
    Use:   "ns",
    Aliases: []string{"NS"},
    Short:  "Resolves NS records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"ns",false)
    },
}

func init() {
    rootCmd.AddCommand(nsCmd)
}