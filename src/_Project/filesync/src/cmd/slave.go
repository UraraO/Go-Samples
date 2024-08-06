/*
 * @Author: chaidaxuan chaidaxuan@wps.cn
 * @Date: 2024-07-31 17:11:30
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-31 18:20:12
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/cmd/slave.go
 * @Description:
 *
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
 */

package cmd

import (
	"filesync/src/def"
	"filesync/src/server"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var slaveCmd = &cobra.Command{
	Use:   "slave",
	Short: "slave run a slave server, with only [Read] function enabled, include Download and List",
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		fmt.Println("files store path is:", filepath)
		encS, err := cmd.Flags().GetString("server-side-encryption")
		if err != nil {
			fmt.Println("server-side-encryption flag error", err)
			return
		}
		encS = strings.ToLower(encS)
		enc := false
		switch encS {
		case "false":
			enc = false
		case "0":
			enc = false
		case "true":
			enc = true
		case "1":
			enc = true
		}
		httpIP, err := cmd.Flags().GetString("http-listen-addr")
		if err != nil {
			fmt.Println("http-listen-addr flag error", err)
			return
		}
		httpPort, err := cmd.Flags().GetInt("http-listen-port")
		if err != nil {
			fmt.Println("http-listen-port flag error", err)
			return
		}
		rpcIP, err := cmd.Flags().GetString("grpc-listen-addr")
		if err != nil {
			fmt.Println("grpc-listen-addr flag error", err)
			return
		}
		rpcPort, err := cmd.Flags().GetInt("grpc-listen-port")
		if err != nil {
			fmt.Println("grpc-listen-port flag error", err)
			return
		}
		server, err := server.NewServer(def.ROLE_SLAVE, httpIP, httpPort, rpcIP, rpcPort, filepath, enc)
		if err != nil {
			fmt.Println("NewServer error:", err.Error())
			return
		}
		fmt.Printf("server-side-encryption: %v, http-listen-addr: %v, http-listen-port: %v, grpc-listen-addr: %v, grpc-listen-port: %v, filepath: %v\n", enc, httpIP, httpPort, rpcIP, rpcPort, filepath)
		server.Run()
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	slaveCmd.Flags().StringP("message", "m", "Hello World", "Message to print")
	slaveCmd.Flags().StringP("server-side-encryption", "", "", "server side encryption, use 'false' or 'true'")
	slaveCmd.Flags().StringP("http-listen-addr", "", "", "http listen addr")
	slaveCmd.Flags().IntP("http-listen-port", "", -1, "http listen port")
	slaveCmd.Flags().StringP("grpc-listen-addr", "", "", "grpc listen addr")
	slaveCmd.Flags().IntP("grpc-listen-port", "", -1, "grpc listen port")
	rootCmd.AddCommand(slaveCmd)
}
