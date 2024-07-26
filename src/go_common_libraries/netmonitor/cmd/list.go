package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list monitoring targets",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		// targets := make([]string, 0, 3)
		// if len(args) > 0 {
		// 	for i := range args {
		// 		targets = append(targets, args[i])
		// 	}
		// }
		initHCM()
		fmt.Println("list targets:")
		HCM.ListAll()
	},
}

func init() {
	// listCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(listCmd)
}
