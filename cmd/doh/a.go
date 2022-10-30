package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var aCmd = &cobra.Command{
    Use:   "a [domain]",
    Aliases: []string{"A"},
    Short:  "Resolves A records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"a",false)
    },
}

func init() {
    rootCmd.AddCommand(aCmd)
}