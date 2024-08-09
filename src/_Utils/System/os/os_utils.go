/*=============
 Author: chaidaxuan chaidaxuan@wps.cn
 Date: 2024-08-08 17:34:06
 LastEditors: chaidaxuan chaidaxuan@wps.cn
 LastEditTime: 2024-08-08 17:34:40
 FilePath: /Golang-Samples/src/_Utils/System/os/os_utils.go
 Description:



 Copyright (c) 2024 by chaidaxuan, All Rights Reserved.
=============*/

package os_utils

import (
	"fmt"
	"os/exec"
	"strings"
)

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
