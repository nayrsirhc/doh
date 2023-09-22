package doh

import (
    "github.com/nayrsirhc/doh/pkg/doh"
    "github.com/spf13/cobra"
)

var anyCmd = &cobra.Command{
    Use:   "any [domain]",
    Aliases: []string{"ANY", "resolve", "get"},
    Short:  "Resolve whatever record type you find for this domain",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        queryName := args[0]
        json, _ := cmd.Flags().GetBool("json")
        doh.RunQuery(queryName,"Not Specified",false, json)
    },
}

func init() {
    rootCmd.AddCommand(anyCmd)
    anyCmd.PersistentFlags().BoolP("json", "j", false, "Return response in JSON Format")
}
