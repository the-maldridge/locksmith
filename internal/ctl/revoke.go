package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var revokeCmd = &cobra.Command{
	Use:     "revoke",
	Short:   "Invalidate a locksmith profile",
	Long:    `Invalidate the supplied locksmith profile.`,
	Example: "telephone revoke <profile>",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("revoke called")
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)
}
