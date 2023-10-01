package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var mxCmd = &cobra.Command{
    Use:   "mx [domain]",
    Aliases: []string{"MX", "mail", "MAIL"},
    Short:  "Resolves MX/MAIL records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"mx",false, json)
    },
}

func init() {
    rootCmd.AddCommand(mxCmd)
    mxCmd.Flags().BoolP("json", "j", false, "Output JSON")
}
