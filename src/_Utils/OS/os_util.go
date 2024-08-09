/*=============
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-08 19:59:16
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-08 19:59:29
 FilePath: /Golang-Samples/src/_Utils/OS/os_util.go
 Description:

 OS Utils

 Copyright (c) 2024 by UraraO, All Rights Reserved.
=============*/

package os_utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// 获取系统进程数
func GetProcessNum() int {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running ps command:", err)
		return -1
	}

	lines := strings.Split(string(output), "\n")
	numProcesses := len(lines) - 1 // Subtract 1 to exclude the header line

	// fmt.Println("Number of processes running:", numProcesses)
	return numProcesses
}
