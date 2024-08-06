/*===========
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:47:01
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-06 22:16:32
 FilePath: /Golang-Samples/src/_Library/cobra/src/cmd/root.go
 Description:

 	cobra命令行框架实例代码
	包含子命令，命令参数，flag的配置

 Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
===========*/

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
