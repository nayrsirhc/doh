package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var aaaaCmd = &cobra.Command{
    Use:   "aaaa",
    Aliases: []string{"AAAA"},
    Short:  "Resolves AAAA records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"aaaa",false)
    },
}

func init() {
    rootCmd.AddCommand(aaaaCmd)
}