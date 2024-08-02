/*
 * @Author: chaidaxuan chaidaxuan@wps.cn
 * @Date: 2024-07-26 18:06:52
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-31 16:56:10
 * @FilePath: /cobra/src/cmd/say.go
 * @Description:
 *
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
 */
package cobratest

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "Print a message",
	Run: func(cmd *cobra.Command, args []string) {
		// message, _ := cmd.Flags().GetString("message")
		// fmt.Println(message)
		fmt.Println(args[0])
	},
	Args: cobra.ExactArgs(1),
}

// root say -m 123
// root say 123
func init() {
	sayCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	rootCmd.AddCommand(sayCmd)
}
