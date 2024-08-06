/*
  - @Author: chaidaxuan chaidaxuan@wps.cn
  - @Date: 2024-08-01 14:31:35
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-08-02 14:45:28
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/server/server_test.go
  - @Description:

server 的各类测试

  - Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package server_test

import (
	"bytes"
	"filesync/src/def"
	"filesync/src/server"
	"filesync/src/utils"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// test pre define
const RoleMaster = def.ROLE_MASTER
const RoleSlave = def.ROLE_SLAVE
const httpIP = ""
const httpPort = 9988
const RPCIP = ""
const RPCPort = 9989
const MasterFilespath = "./files_master"
const SlaveFilespath = "./files_slave"
const encryption = false

func TestGetFilenameFromURL(t *testing.T) {
	URL := "GET /files/filenamez"

	fname, err := server.GetFileNameFromURL(URL)
	if err != nil {
		t.Error(err)
	}
	t.Log(fname)
}

func TestGetFileList(t *testing.T) {
	path := "./files_master/"
	fs := []string{"test1", "test2", "test3"}
	files := server.GetFilelist(fs, path)
	for i := range files.Files {
		if files.Files[i].FileName != fs[i] {
			t.Error(files.Files[i].FileName, " != ", fs[i])
			t.Fail()
		}
	}
}

// go test -v -run ^TestCheckFileExistorIsDir$ ./src/server
func TestCheckFileExistorIsDir(t *testing.T) {
	fe := "../../logs/server-master.log"
	fne := "../../logs/server-master--.log"
	fd := "../../src/server"
	fmt.Println(os.Getwd())
	fmt.Println(utils.CheckFileExistorIsDir(fe))
	fmt.Println(utils.CheckFileExistorIsDir(fne))
	fmt.Println(utils.CheckFileExistorIsDir(fd))
}

func TestServerList(t *testing.T) {
	s, err := server.NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Files.FilesMap["test1"] = server.NewFile("test1", s.Files)
	s.Files.FilesMap["test2"] = server.NewFile("test2", s.Files)
	s.Files.FilesMap["test3"] = server.NewFile("test3", s.Files)

	go s.Run()
	resp, err := http.Get("http://127.0.0.1:9988/files")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer resp.Body.Close()
	buffer := make([]byte, 4096)
	resp.Body.Read(buffer)
	t.Log(buffer)
}

func TestServerUpload(t *testing.T) {
	s, err := server.NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.Run()
	r := bytes.NewReader([]byte("test content"))
	req, _ := http.NewRequest("PUT", "http://127.0.0.1:9988/files/file1", r)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer resp.Body.Close()
	buffer := make([]byte, 4096)
	resp.Body.Read(buffer)
	t.Log(buffer)
}

func TestServerDownload(t *testing.T) {
	s, err := server.NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.Run()
	resp, err := http.Get("http://127.0.0.1:9988/files/file1")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer resp.Body.Close()
	buffer := make([]byte, 4096)
	resp.Body.Read(buffer)
	t.Log(buffer)
}

func TestServerDelete(t *testing.T) {
	s, err := server.NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.Run()
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:9988/files/file1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(resp.StatusCode)
}

func TestServerSlave(t *testing.T) {
	s, err := server.NewServer(RoleSlave, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Run()
}

func TestServerSyncOneSlave(t *testing.T) {
	s, err := server.NewServer(RoleSlave, httpIP, httpPort, RPCIP, RPCPort, SlaveFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Run()
}

func TestServerSyncOneMaster(t *testing.T) {
	s, err := server.NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.Run()

	time.Sleep(time.Second)
	fmt.Printf("start to sync one upload\n")
	filenameWithoutPrefix := "FILENAMEU"
	err = s.SyncCli.SyncOneUpload(filenameWithoutPrefix)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second * 5)
	fmt.Printf("start to sync one delete\n")
	err = s.SyncCli.SyncOneDelete(filenameWithoutPrefix)
	if err != nil {
		fmt.Println(err)
	}
}

func TestServerSyncOneSlaveFull(t *testing.T) {
	s, err := server.NewServer(RoleSlave, httpIP, httpPort, RPCIP, RPCPort, SlaveFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.RunWithoutHttp()
}

func TestServerSyncOneMasterFull(t *testing.T) {
	s, err := server.NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Run()
}

// go test -benchmem -run=^$ -bench ^BenchmarkServer$ ./src/server
func BenchmarkServer(b *testing.B) {
	M := 100
	N := 100
	wg := sync.WaitGroup{}
	wg.Add(M * N)
	bUploadDownloadListDelete := func(prefix string) {
		for i := 0; i < N; i++ {
			filename := "BenchTestFile" + prefix + strconv.Itoa(i)
			upload, _ := http.NewRequest("PUT", "http://127.0.0.1:9988/files/"+filename, bytes.NewReader([]byte("test content")))
			download, _ := http.NewRequest("GET", "http://127.0.0.1:9988/files/"+filename, nil)
			delete, _ := http.NewRequest("DELETE", "http://127.0.0.1:9988/files/"+filename, nil)
			list, _ := http.NewRequest("GET", "http://127.0.0.1:9988/files?limit=10&offset=0", nil)
			cli := http.DefaultClient
			cli.Do(upload)
			// time.Sleep(30 * time.Millisecond)
			cli.Do(list)
			r1, err := cli.Do(download)
			if err != nil {
				fmt.Println(err)
			} else {
				buffer := make([]byte, 16)
				r1.Body.Read(buffer)
				fmt.Println(string(buffer))
				r1.Body.Close()
			}
			// time.Sleep(30 * time.Millisecond)
			r2, err := cli.Do(delete)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(r2.StatusCode)
			}
			wg.Done()
		}
	}
	beg := time.Now()
	for i := 0; i < M; i++ {
		go bUploadDownloadListDelete(strconv.Itoa(i) + "-")
	}
	wg.Wait()
	end := time.Now()
	fmt.Println("time cost:", end.Sub(beg).Milliseconds())
}

// go test -benchmem -run=^$ -bench ^BenchmarkOutOfOrder$ ./src/server
func BenchmarkOutOfOrder(b *testing.B) {
	M := 10
	N := 100
	wg := sync.WaitGroup{}
	wg.Add(M * N)
	bUploadDownloadListDelete := func(prefix string) {
		for i := 0; i < N; i++ {
			// filename := "BenchTestFile" + prefix + strconv.Itoa(i)
			filename := "BenchTestFile" + prefix
			upload, _ := http.NewRequest("PUT", "http://127.0.0.1:9988/files/"+filename, bytes.NewReader([]byte("test content")))
			download, _ := http.NewRequest("GET", "http://127.0.0.1:9988/files/"+filename, nil)
			delete, _ := http.NewRequest("DELETE", "http://127.0.0.1:9988/files/"+filename, nil)
			// list, _ := http.NewRequest("GET", "http://127.0.0.1:9988/files?limit=10&offset=0", nil)
			cli := http.DefaultClient
			r1, err := cli.Do(download)
			if err != nil {
				fmt.Println("download 1 err:", err)
			} else {
				if r1.StatusCode == http.StatusNotFound {
					fmt.Println("Not Found")
				} else {
					buffer := make([]byte, 16)
					r1.Body.Read(buffer)
					fmt.Println(string(buffer))
					r1.Body.Close()
				}
			}
			cli.Do(upload)
			// time.Sleep(30 * time.Millisecond)
			// cli.Do(list)
			r1, err = cli.Do(download)
			if err != nil {
				fmt.Println("download 2 err:", err)
			} else {
				if r1.StatusCode == http.StatusNotFound {
					fmt.Println("Not Found")
				} else {
					buffer := make([]byte, 16)
					r1.Body.Read(buffer)
					fmt.Println(string(buffer))
					r1.Body.Close()
				}
			}
			r1, err = cli.Do(download)
			if err != nil {
				fmt.Println("download 3 err:", err)
			} else {
				if r1.StatusCode == http.StatusNotFound {
					fmt.Println("Not Found")
				} else {
					buffer := make([]byte, 16)
					r1.Body.Read(buffer)
					fmt.Println(string(buffer))
					r1.Body.Close()
				}
			}
			// time.Sleep(30 * time.Millisecond)
			r2, err := cli.Do(delete)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(r2.StatusCode)
			}
			r1, err = cli.Do(download)
			if err != nil {
				fmt.Println("download 4 err:", err)
			} else {
				if r1.StatusCode == http.StatusNotFound {
					fmt.Println("Not Found")
				} else {
					buffer := make([]byte, 16)
					r1.Body.Read(buffer)
					fmt.Println(string(buffer))
					r1.Body.Close()
				}
			}
			wg.Done()
		}
	}
	time.Sleep(time.Second)
	beg := time.Now()
	for i := 0; i < M; i++ {
		go bUploadDownloadListDelete(strconv.Itoa(i) + "-")
	}
	wg.Wait()
	end := time.Now()
	fmt.Println("time cost:", end.Sub(beg).Milliseconds())
}
