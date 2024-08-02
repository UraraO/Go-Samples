/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-29 14:09:10
  - @LastEditors: chaidaxuan chaidaxuan@wps.cn
  - @LastEditTime: 2024-08-02 11:04:56
  - @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/server/server.go

* @Description:

filesync文件服务器的类型定义，及基本接口
包括创建Server对象，以及阻塞的Run方法
一个filesync服务器的服务主要分3部分：Http，File，Sync

* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package server

import (
	"filesync/src/def"
	"filesync/src/logger"
	"filesync/src/utils"
	"fmt"
	"os"
	"sync"
)

// 通用IP地址结构
type ServerAddr struct {
	IP   string
	Port int
}

func (sa *ServerAddr) Address() string {
	return fmt.Sprintf("%s:%d", sa.IP, sa.Port)
}

// 文件服务器类，同时作为Master和Slave
type Server struct {
	// server role config
	Role       int
	HttpAddr   ServerAddr
	MasterAddr ServerAddr // 作为Slave时，该项为监听的RPC地址
	SlaveAddr  ServerAddr // 作为Master时，该项为Slave的RPC地址
	// 上方冗余设计本意是想实现双向通信，便于双向的全量同步，当前仅使用Master请求，Slave响应

	// sync 同步服务部分，当前Slave作为RPC服务端，Master请求同步单个文件和保活
	SyncServ     SyncServer
	SyncCli      SyncClient
	FullSyncMut  sync.RWMutex
	IsInFullSync bool // 是否正在全量同步，需要锁双方的全局文件列表
	// sync queue，file协程放入任务，sync协程取出任务
	// SyncQueue *utils.BlockQueue

	// mission queue，http协程放入任务，file协程取出任务
	// MissionQueue *utils.BlockQueue

	// file serve
	// FilelistSnapshot *FileList
	Files *FileList // 全部文件列表
	// FilesInUse    *FileList // 正在操作的文件
	FilesPath     string // 文件存储路径
	UseEncryption bool   // 服务端加密

	// log
	Logger *logger.ULog
}

var IsMaster = true

// 新建服务器，指定服务器角色（M/S)
func NewServer(role int, httpIP string, httpPort int, RPCIP string, RPCPort int, filespath string, encryption bool) (server *Server, err error) {
	if httpIP == "127.0.0.1" || httpIP == "localhost" {
		httpIP = ""
	}
	if RPCIP == "127.0.0.1" || RPCIP == "localhost" {
		RPCIP = ""
	}
	server = &Server{
		Role: def.ROLE_INIT,
		HttpAddr: ServerAddr{
			IP:   httpIP,
			Port: httpPort,
		},
		FullSyncMut:  sync.RWMutex{},
		IsInFullSync: false,
		// SyncQueue:     utils.NewBlockQueue(100),
		// MissionQueue:  utils.NewBlockQueue(1000),
		Files: NewFileList(),
		// FilesInUse:    NewFileList(),
		FilesPath:     filespath + "/",
		UseEncryption: encryption,
		// Logger:        *logger.NewLogger("server-master.log"),
	}
	if err := utils.CreateDir(filespath); err != nil {
		return nil, err
	}
	IsMaster = (role == def.ROLE_MASTER)
	if role == def.ROLE_MASTER {
		server.Logger = logger.NewLogger(def.DEFAULT_LOG_PATH + def.DEFAULT_LOG_NAME_MASTER)
		server.MasterAddr = ServerAddr{
			IP:   httpIP,
			Port: httpPort,
		}
		server.SlaveAddr = ServerAddr{
			IP:   RPCIP,
			Port: RPCPort,
		}
	} else if role == def.ROLE_SLAVE {
		server.Logger = logger.NewLogger(def.DEFAULT_LOG_PATH + def.DEFAULT_LOG_NAME_SLAVE)
		server.SlaveAddr = ServerAddr{
			IP:   httpIP,
			Port: httpPort,
		}
		server.MasterAddr = ServerAddr{
			IP:   RPCIP,
			Port: RPCPort,
		}
	} else {
		return nil, fmt.Errorf("unknown server role or role[INIT], please use ROLE_MASTER or ROLE_SLAVE")
	}
	server.SyncServ = *NewSyncServer(server)
	server.SyncCli = *NewSyncClient(server)
	server.Role = role

	err = LoadFilesInDir(server)
	if err != nil {
		server.Logger.Printf("server init: LoadFilesInDir error: %v\n", err)
		return nil, err
	}
	fmt.Println("server init ok")
	server.PrintServerFiles()
	return server, nil
}

func (s *Server) PrintServerFiles() {
	fmt.Println("*** Files:")
	for fname, f := range s.Files.FilesMap {
		fmt.Println("*** ", fname, f.Info.Size, f.Info.CTime, f.Info.MTime)
	}
}

func LoadFilesInDir(server *Server) error {
	flist, err := os.ReadDir(server.FilesPath)
	if err != nil {
		server.Logger.Printf("server init: os.ReadDir error: %v\n", err)
		return err
	}
	server.Files.Mut.Lock()
	defer server.Files.Mut.Unlock()
	for _, f := range flist {
		if f.IsDir() {
			continue
		}
		server.Files.FilesMap[f.Name()] = NewFile(f.Name(), server.Files)
		info, err := f.Info()
		if err != nil {
			server.Logger.Printf("server init: f.Info() error: %v\n", err)
			continue
		}
		server.Files.FilesMap[f.Name()].Info = &FileInfo{
			FileName: f.Name(),
			Size:     info.Size(),
			CTime:    info.ModTime().UnixMilli(),
			MTime:    info.ModTime().UnixMilli(),
		}
	}
	return nil
}

func (s *Server) Run() error {
	// 启动http服务
	// go 文件服务，消费任务

	// RPC part
	// if master
	// go RPC客户端，向slave发送sync任务

	// if slave
	// go 启动RPC服务，接收sync任务

	go s.RunFileServer()
	go s.RunSyncServer()
	return s.RunHttpServer()
}

func (s *Server) RunWithoutHttp() error {
	// 启动http服务
	// go 文件服务，消费任务

	// RPC part
	// if master
	// go RPC客户端，向slave发送sync任务

	// if slave
	// go 启动RPC服务，接收sync任务

	go s.RunFileServer()
	s.RunSyncServer()
	return nil
	// return s.RunHttpServer()
}

func (s *Server) RunWithoutSync() error {
	// 启动http服务
	// go 文件服务，消费任务

	// RPC part
	// if master
	// go RPC客户端，向slave发送sync任务

	// if slave
	// go 启动RPC服务，接收sync任务

	go s.RunFileServer()
	// s.RunSyncServer()
	return s.RunHttpServer()
}
