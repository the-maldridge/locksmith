package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:     "approve",
	Short:   "Approve validates a locksmith profile",
	Long:    `Approve is used to validate a profile for use with locksmith`,
	Example: "telephone approve <profile>",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("approve called")
	},
}

func init() {
	rootCmd.AddCommand(approveCmd)
}
