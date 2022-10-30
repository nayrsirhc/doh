package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var anyCmd = &cobra.Command{
    Use:   "any [domain]",
    Aliases: []string{"ANY", "resolve", "get"},
    Short:  "Resolve whatever record type you find for this domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"Not Specified",false)
    },
}

func init() {
    rootCmd.AddCommand(anyCmd)
}