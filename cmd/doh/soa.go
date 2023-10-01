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
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"soa",false, json)
    },
}

func init() {
    rootCmd.AddCommand(soaCmd)
    soaCmd.Flags().BoolP("json", "j", false, "Output JSON")
}
