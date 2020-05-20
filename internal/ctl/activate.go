package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var activateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Enable a locksmith profile",
	Long:    `Enable the supplied locksmith profile.`,
	Example: "telephone activate <profile>",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("activate called")
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
