package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var nsCmd = &cobra.Command{
    Use:   "ns [domain]",
    Aliases: []string{"NS"},
    Short:  "Resolves NS records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"ns",false, json)
    },
}

func init() {
    rootCmd.AddCommand(nsCmd)
    nsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
