/*
 * @Author: chaidaxuan chaidaxuan@wps.cn
 * @Date: 2024-07-31 16:33:11
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-31 17:46:18
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/cmd/root.go
 * @Description:
 *
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
 */
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync is a file server with ,aster-slave sync function",
	Long:  `sync is a file server with ,aster-slave sync function`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd, args, errors.New("unrecognized command"))
	},
}

func Execute() {
	rootCmd.Execute()
}
