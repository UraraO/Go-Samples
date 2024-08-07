/*=============
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:47:01
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-07 22:05:51
 FilePath: /Golang-Samples/src/_Utils/Net/rpc/rpcServer/server.go
 Description:

 RPC server by golang.org

 Copyright (c) 2024 by UraraO, All Rights Reserved.
=============*/

package main

import (
	rpccommon "Golang-Samples/src/_Utils/Net/rpc/rpcCommon"
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
