/*=============
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:47:01
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-07 22:06:09
 FilePath: /Golang-Samples/src/_Utils/Net/rpc/rpcClient/client.go
 Description:

 RPC client by golang.org

 Copyright (c) 2024 by UraraO, All Rights Reserved.
=============*/

package main

import (
	rpccommon "Golang-Samples/src/_Utils/Net/rpc/rpcCommon"
	"fmt"
	"net/rpc"
)

func ClientTest() {
	cli, err := rpc.DialHTTP("tcp", ":9988")
	if err != nil {
		fmt.Println(err)
		return
	}
	sum := 0
	cli.Call("numServer.Add", &rpccommon.Args{1, 2}, &sum)
	fmt.Println(sum)
}

func main() {
	ClientTest()
}
