package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var soaCmd = &cobra.Command{
    Use:   "soa [domain]",
    Aliases: []string{"SOA"},
    Short:  "Resolves SOA for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"soa",false)
    },
}

func init() {
    rootCmd.AddCommand(soaCmd)
}