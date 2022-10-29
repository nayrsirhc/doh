package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var txtCmd = &cobra.Command{
    Use:   "txt",
    Aliases: []string{"TXT"},
    Short:  "Resolves TXT records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"txt",false)
    },
}

func init() {
    rootCmd.AddCommand(txtCmd)
}