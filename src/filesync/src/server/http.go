/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-29 16:06:28
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-08-01 16:56:20
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/server/http.go

* @Description:

文件服务器的Http服务，对外提供4个http接口：Upload，Download，List，Delete

* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package server

import (
	"filesync/src/def"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) RunHttpServer() error {
	// 监听http端口，启动http服务
	// http服务生产任务，写入s.MissionQueue
	// file服务取出任务，消费

	// master / slave
	// read: List / Download
	http.HandleFunc("GET /files/{filename}", s.HandleDownload)
	http.HandleFunc("GET /files", s.HandleList) // List入口
	// master
	// write: Upload / Delete
	if IsMaster {
		http.HandleFunc("PUT /files/{filename}", s.HandleUpload)
		http.HandleFunc("DELETE /files/{filename}", s.HandleDelete)
	}

	server := &http.Server{Addr: s.HttpAddr.Address(), Handler: nil}
	ln, err := net.Listen("tcp", s.HttpAddr.Address())
	if err != nil {
		return err
	}

	return server.Serve(ln)
}

// http://localhost:9988/files/FILENAMEUPLOAD
func (s *Server) HandleUpload(w http.ResponseWriter, r *http.Request) {
	s.Logger.Printf("[Upload] New Query: %s\n", r.URL.RequestURI())
	beg := time.Now()
	if !IsMaster {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	filename := r.PathValue("filename")
	// 上传文件
	// TODO，利用MissionQueue，向FileServer提交任务
	err := s.UploadFile(filename, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fmt.Println("upload success, now Files is:", s.Files.FilesMap)
	// s.PrintServerFiles()
	w.WriteHeader(http.StatusOK)
	s.Logger.Printf("delete success, filename: %s, cost time: %v ms\n", filename, time.Since(beg).Milliseconds())
}

// http://localhost:9988/files/FILENAMEDOWNLOAD
func (s *Server) HandleDownload(w http.ResponseWriter, r *http.Request) {
	s.Logger.Printf("[Download] New Query: %s\n", r.URL.RequestURI())
	beg := time.Now()
	filenameWithoutPrefix := r.PathValue("filename")
	// w.Write([]byte(fmt.Sprintf("fileName: %s", filename)))

	// 下载文件
	// TODO，利用MissionQueue，向FileServer提交任务
	err := s.DownloadFile(w, r, filenameWithoutPrefix)
	if err != nil {
		s.Logger.Println("download error: ", err)
		return
	}

	// s.Logger.Println("download success, filename: ", filenameWithoutPrefix)
	s.Logger.Printf("download success, filename: %s, cost time: %v ms\n", filenameWithoutPrefix, time.Since(beg).Milliseconds())
}

// http://localhost:9988/files/FILENAMEDELETE
func (s *Server) HandleDelete(w http.ResponseWriter, r *http.Request) {
	s.Logger.Printf("[Delete] New Query: %s\n", r.URL.RequestURI())
	beg := time.Now()
	if !IsMaster {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	filename := r.PathValue("filename")
	// 删除文件
	// TODO，利用MissionQueue，向FileServer提交任务
	err := s.DeleteFile(filename)
	if err != nil {
		if err.Error() != def.FILE_NOT_EXIST_IN_FILES && err.Error() != def.FILE_NOT_EXIST && err.Error() != def.FILE_IS_DIR {
			w.WriteHeader(http.StatusNoContent)
			s.Logger.Println("download error: ", err)
			return
		}
	}
	// fmt.Println("delete success, filename:", filename)
	// s.PrintServerFiles()
	w.WriteHeader(http.StatusNoContent)
	s.Logger.Printf("delete success, filename: %s, cost time: %v ms\n", filename, time.Since(beg).Milliseconds())
}

// http://localhost:9988/files?limit=5&offset=10
func (s *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	s.Logger.Printf("[List] New Query: %s\n", r.URL.RequestURI())
	beg := time.Now()
	// 接收参数
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		fmt.Println("limit err: ", err)
	}
	if limit > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// fmt.Printf("limit: %v \n", limit)
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		fmt.Println("offset err: ", err)
	}
	// fmt.Printf("offset: %v \n", offset)

	// 获取文件列表JSON数据
	// TODO，利用MissionQueue，向FileServer提交任务
	b := s.GetFilelistJSON(limit, offset)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	s.Logger.Printf("list success, cost time: %v ms\n", time.Since(beg).Milliseconds())
}

// ***********************************
// http://localhost:9988/upload/fileuuuu
func HttpTemplate(addr string, handler http.Handler) error {
	http.HandleFunc("/", getHandle)
	http.HandleFunc("/upload/", uploadHandler)
	server := &http.Server{Addr: addr, Handler: nil}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return server.Serve(ln)
}
