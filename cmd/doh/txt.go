package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var txtCmd = &cobra.Command{
    Use:   "txt [domain]",
    Aliases: []string{"TXT"},
    Short:  "Resolves TXT records for domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"txt",false, json)
    },
}

func init() {
    rootCmd.AddCommand(txtCmd)
    txtCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
