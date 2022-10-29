package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var extensiveCmd = &cobra.Command{
    Use:   "extensive",
    Aliases: []string{"ext"},
    Short:  "Resolves an extensive list of records for a domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.QueryExtensive(queryName)
    },
}

func init() {
    rootCmd.AddCommand(extensiveCmd)
}