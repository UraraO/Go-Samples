package main

import (
	rpccommon "Golang-Samples/src/Net/rpc/rpcCommon"
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
