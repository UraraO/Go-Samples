package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new monitoring target(such as domain name)",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		initHCM()
		targets := make([]string, 0, 3)
		if len(args) > 0 {
			targets = append(targets, args...)
		}
		fmt.Println("add targets:", targets)
		for i := range targets {
			HCM.AddHC(dur, targets[i], handler)
		}
	},
}

func init() {
	// addCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(addCmd)
}
