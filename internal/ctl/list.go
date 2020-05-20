package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all locksmith profiles",
	Long: `Display all of the locksmith profiles configured on the host
machine`,
	Example: "telephone profile list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	profileCmd.AddCommand(listCmd)
}
