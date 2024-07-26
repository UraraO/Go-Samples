package cobratest

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print go version",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		// fmt.Println()
		govers := exec.Command("go", "version")
		out, _ := govers.Output()
		fmt.Println(string(out))
	},
}

func init() {
	// versionCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(versionCmd)
}
