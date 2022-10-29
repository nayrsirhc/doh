package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var mxCmd = &cobra.Command{
    Use:   "mx",
    Aliases: []string{"MX", "mail", "MAIL"},
    Short:  "Resolves MX/MAIL records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"mx",false)
    },
}

func init() {
    rootCmd.AddCommand(mxCmd)
}