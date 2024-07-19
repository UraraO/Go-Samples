package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop a monitoring target or stop all targets",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		// targets := make([]string, 0, 3)
		// if len(args) > 0 {
		// 	for i := range args {
		// 		targets = append(targets, args[i])
		// 	}
		// }
		initHCM()
		fmt.Println("stop targets:")
		HCM.StopAll()
	},
}

func init() {
	// addCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(stopCmd)
}
