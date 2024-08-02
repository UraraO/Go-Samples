/*
 * @Author: chaidaxuan chaidaxuan@wps.cn
 * @Date: 2024-07-26 18:06:52
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-31 16:46:14
 * @FilePath: /cobra/src/cmd/env.go
 * @Description:
 *
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
 */
package cobratest

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print environment message",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		var goenv *exec.Cmd
		if len(name) == 0 {
			goenv = exec.Command("go", "env")
		} else {
			goenv = exec.Command("go", "env", name)
		}
		out, _ := goenv.Output()
		fmt.Println(string(out))
	},
}

func init() {
	envCmd.Flags().StringP("name", "n", "", "env to print")
	rootCmd.AddCommand(envCmd)
}
