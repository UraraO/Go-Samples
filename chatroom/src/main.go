package main

import serve "IM_system/src/server"

func main() {
	serve.MainClient()
	// serve.MainServer()
	select {}
}
