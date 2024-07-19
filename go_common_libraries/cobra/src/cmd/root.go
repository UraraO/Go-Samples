package cobratest

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "roor is a distributed version control system.",
	Long: `root is a free and open source distributed version control system
  designed to handle everything from small to very large projects 
  with speed and efficiency.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd, args, errors.New("unrecognized command"))
	},
}

func Execute() {
	rootCmd.Execute()
}
