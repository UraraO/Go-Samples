package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a monitoring target",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		// targets := make([]string, 0, 3)
		// if len(args) > 0 {
		// 	for i := range args {
		// 		targets = append(targets, args[i])
		// 	}
		// }
		// fmt.Println("")
		initHCM()
		HCM.StartAll()
	},
}

func init() {
	// addCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(startCmd)
}
