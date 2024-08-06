/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-30 17:12:42
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-30 17:29:25
 * @FilePath: /urarao/GoProjects/Golang-Samples/src/Net/rpc/rpcServer/server.go

* @Description:

*
* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package main

import (
	rpccommon "Golang-Samples/src/Net/rpc/rpcCommon"
	"net/http"
	"net/rpc"
)

func RPCServerTest() {
	server := &rpccommon.RPCServer{}
	rpc.Register(server)
	rpc.RegisterName("numServer", server)
	rpc.HandleHTTP()
	http.ListenAndServe(":9988", nil)
}

func main() {
	RPCServerTest()
}
