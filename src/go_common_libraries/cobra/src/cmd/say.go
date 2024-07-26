package cobratest

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "Print a message",
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		fmt.Println(message)
	},
}

func init() {
	sayCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(sayCmd)
}
