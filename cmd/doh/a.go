package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var aCmd = &cobra.Command{
    Use:   "a",
    Short:  "Reverses a string",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        doh.RunQuery(queryName,"A",false)
    },
}

func init() {
    rootCmd.AddCommand(aCmd)
}