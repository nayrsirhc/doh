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
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"Not Specified",false, json)
    },
}

func init() {
    rootCmd.AddCommand(cnameCmd)
    cnameCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
