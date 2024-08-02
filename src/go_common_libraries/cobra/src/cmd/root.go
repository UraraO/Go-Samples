/*
 * @Author: chaidaxuan chaidaxuan@wps.cn
 * @Date: 2024-07-26 18:06:52
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-31 16:40:58
 * @FilePath: /cobra/src/cmd/root.go
 * @Description:
 *
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
 */
package cobratest

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "root is a distributed version control system.",
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
