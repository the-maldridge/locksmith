package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate a locksmith profile",
	Long:    `Deactivate removes the supplied locksmith profile from use`,
	Example: "telephone deactivate <profile>",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deactivate called")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
