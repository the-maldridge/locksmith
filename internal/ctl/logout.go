package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "Logout of a locksmith profile",
	Long:    `Logout of the supplied locksmith profile `,
	Example: "telephone logout <profile>",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logout called")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
