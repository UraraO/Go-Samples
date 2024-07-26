package cmd

import (
	"errors"
	"fmt"
	"netmonitor/src/heartcheck"
	"time"

	"github.com/spf13/cobra"
)

var HCM *heartcheck.HCManager
var dur time.Duration = 5 * time.Second
var handler = func() {
	fmt.Println("handler handling------")
}

var rootCmd = &cobra.Command{
	Use:   "netmonitor",
	Short: "nm",
	Long:  `netmonitor is a net monitor, which provide net heart ticker tool, and manage all tickers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd, args, errors.New("unrecognized command"))
	},
}

func initHCM() {
	if HCM == nil {
		HCM = heartcheck.InitHCManager()
	}
}

func Execute() {
	rootCmd.Execute()
}
