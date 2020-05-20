package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   "Create a locksmith profile",
	Long:    `Interactively generate a locksmith profile configuration`,
	Example: "telephone profile setup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
	},
}

func init() {
	profileCmd.AddCommand(setupCmd)
}
