package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var aaaaCmd = &cobra.Command{
    Use:   "aaaa [domain]",
    Aliases: []string{"AAAA"},
    Short:  "Resolves AAAA records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"aaaa",false, json)
    },
}

func init() {
    rootCmd.AddCommand(aaaaCmd)
    aaaaCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
