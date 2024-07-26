package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "make report of monitoring targets",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		// targets := make([]string, 0, 3)
		// if len(args) > 0 {
		// 	for i := range args {
		// 		targets = append(targets, args[i])
		// 	}
		// }
		initHCM()
		fmt.Println("made report of targets:")
		HCM.ReportAll()
	},
}

func init() {
	// addCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(reportCmd)
}
