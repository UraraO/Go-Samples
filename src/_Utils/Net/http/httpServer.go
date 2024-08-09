/*=============
 Author: chaidaxuan chaidaxuan@wps.cn
 Date: 2024-08-07 15:07:50
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-09 20:29:27
 FilePath: /Golang-Samples/src/_Utils/Net/http/httpServer.go
 Description:



 Copyright (c) 2024 by chaidaxuan, All Rights Reserved.
=============*/

package http_test

import (
	file_utils "Golang-Samples/src/_Utils/File"
	"fmt"
	"net"
	"net/http"
)

type ServerAddr struct {
	IP   string
	Port int
}

func (sa *ServerAddr) Address() string {
	return fmt.Sprintf("%s:%d", sa.IP, sa.Port)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello\n"))
}

func handleNull(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to UraraO's HTTP Server\n"))
}

// GET http://localhost:9988/files/filename
// 将服务器的文件下载给客户端
func handleFileDownload(w http.ResponseWriter, r *http.Request) {
	filenameWithoutPrefix := r.PathValue("filename")
	filepath := "./" + filenameWithoutPrefix
	exist, isdir, info := file_utils.CheckFileExistorIsDir(filepath)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
	} else if exist && isdir {
		w.WriteHeader(http.StatusForbidden)
	}
	w.Header().Set("Last-Modified", info.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	// 该函数可以自动分块的将一个文件发送给http请求方
	http.ServeFile(w, r, filepath)
}

func SimpleHttpServerTest(ip string, port int) error {
	addr := ServerAddr{
		IP:   ip,
		Port: port,
	}
	server := &http.Server{Addr: addr.Address(), Handler: nil}
	ln, err := net.Listen("tcp", addr.Address())
	if err != nil {
		return err
	}

	http.HandleFunc("GET /hello", handleHello)
	http.HandleFunc("GET /", handleNull)
	http.HandleFunc("GET /files/{filename}", handleFileDownload)

	return server.Serve(ln)
}
