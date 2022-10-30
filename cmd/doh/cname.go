package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var cnameCmd = &cobra.Command{
    Use:   "cname [domain]",
    Aliases: []string{"CNAME"},
    Short:  "Resolves CNAME records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"cname",false)
    },
}

func init() {
    rootCmd.AddCommand(cnameCmd)
}