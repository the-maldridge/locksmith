package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Short:   "Profile management",
	Long:    `Manage locksmith profiles`,
	Example: "telephone profile <subcommand>",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile called")
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
