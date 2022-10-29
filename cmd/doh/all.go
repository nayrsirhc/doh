package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
    Use:   "all",
    Short:  "Reverses a string",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.QueryAll(queryName)
    },
}

func init() {
    rootCmd.AddCommand(allCmd)
}