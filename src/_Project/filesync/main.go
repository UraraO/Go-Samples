/*
  - @Author: chaidaxuan chaidaxuan@wps.cn
  - @Date: 2024-07-29 09:34:39
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-08-01 16:41:24
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/main.go
  - @Description:

filesync文件服务器
编译为sync可执行文件
// 运行命令：
// 启动master服务器
// ./sync master ./files_master --http-listen-addr= --http-listen-port=9988 --slave-grpc-addr= --slave-grpc-port=9989 --server-side-encryption false
// 启动slave服务器
// ./sync slave ./files_slave --http-listen-addr= --http-listen-port=9990 --grpc-listen-addr= --grpc-listen-port=9989 --server-side-encryption false

  - Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package main

import "filesync/src/cmd"

func main() {
	cmd.Execute()
}
