/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-29 17:00:42
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-08-01 16:55:44
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/server/sync.go

* @Description:

文件服务器的主从同步方法，利用RPC，实现了在RPC服务端提供同步接口，客户端申请将文件同步到服务端
实际使用中，从机为RPC服务端，主机接到文件写请求时，申请将文件同步到从机

* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package server

import (
	"bytes"
	"errors"
	"filesync/src/def"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

const TestSyncServerAddr = ":9988"
const MaxSyncOneRetry = 10

type SyncServer struct {
	From       *Server
	ServerName string
}

func NewSyncServer(from *Server) *SyncServer {
	server := &SyncServer{
		ServerName: "SyncServer",
		From:       from,
	}
	return server

}

func (s *SyncServer) Run(Addr string) {
	rpc.Register(s)
	rpc.RegisterName(s.ServerName, s)
	rpc.HandleHTTP()
	http.ListenAndServe(Addr, nil)
}

type SyncClient struct {
	From       *Server
	ClientName string
	cli        *rpc.Client
	Connected  bool
}

func NewSyncClient(from *Server) *SyncClient {
	cli := &SyncClient{
		ClientName: "SyncClient",
		From:       from,
	}
	return cli
}

func (c *SyncClient) Connect(Addr string) error {
	if !c.Connected {
		cli, err := rpc.DialHTTP("tcp", Addr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		c.cli = cli
		c.Connected = true
	}
	return nil
}

// 确认存活方法
type IsAliveArgs struct {
}

type IsAliveReply struct {
	OK bool
}

func (s *SyncServer) IsAlive(args IsAliveArgs, reply *IsAliveReply) error {
	*reply = IsAliveReply{true}
	return nil
}

func (c *SyncClient) SyncIsAlive() bool {
	if !c.Connected {
		err := c.Connect(c.From.SlaveAddr.Address())
		if err != nil {
			c.From.Logger.Printf("SyncIsAlive try to connect error, %v\n", err.Error())
			return false
		}
	}

	rep := IsAliveReply{false}
	err := c.cli.Call("SyncServer.IsAlive", &IsAliveArgs{}, &rep)
	if err != nil {
		c.From.Logger.Printf("RPC cli Call SyncIsAlive error, %v\n", err.Error())
		return false
	}
	return rep.OK
}

func (c *SyncClient) RunAliveCheck() {
	alive, lastalive := false, true
	for {
		alive = c.SyncIsAlive()
		if alive {
			// c.From.Logger.Printf("Slave server is alive, slave: %v\n", c.From.SlaveAddr.Address())
			if !lastalive { // 重新上线
				// 全量同步
				// c.SyncAll()
				c.From.Logger.Printf("Slave server is re-on-line, try to SyncAll, slave: %v\n", c.From.SlaveAddr.Address())
				ok := c.SyncAll()
				if ok {
					c.From.Logger.Printf("SyncAll OK, slave: %v\n", c.From.SlaveAddr.Address())
				} else {
					c.From.Logger.Printf("SyncAll has some error, slave: %v\n", c.From.SlaveAddr.Address())
				}
			}
		} else {
			c.From.Logger.Printf("Slave server is off-line, try to reconnect, slave: %v\n", c.From.SlaveAddr.Address())
			c.Connected = false
			if c.cli != nil {
				c.cli.Close()
			}
		}
		lastalive = alive
		time.Sleep(time.Second)
	}
}

// 同步一个文件
type SyncOneArgs struct {
	Operation   string
	Filename    string
	NeedContent bool
	FileContent []byte

	// if big
	IsBig    bool
	ChunkNum int64
	First    bool
	Last     bool
	Which    int64
}

type SyncOneReply struct {
	OK bool
}

func (s *SyncServer) SyncOne(args SyncOneArgs, reply *SyncOneReply) error {
	// fmt.Println("SyncOne: ", args.Operation, args.Filename, args.NeedContent, string(args.FileContent))
	s.From.Logger.Println("SyncOne: ", args.Operation, args.Filename, args.NeedContent, args.IsBig, args.ChunkNum, args.First, args.Last)
	var err error
	switch args.Operation {
	case def.OP_UPLOAD:
		err = s.SyncOneFileUpload(args.Filename, args.FileContent)
	case def.OP_DELETE:
		err = s.SyncOneFileDelete(args.Filename)
	case def.OP_UPLOAD_BIG:
		err = s.SyncOneBigUpload(args.Filename, args.FileContent, args.First, args.Last)
	}
	if err != nil {
		*reply = SyncOneReply{false}
		return err
	}
	*reply = SyncOneReply{true}
	return nil
}

func (s *SyncServer) SyncOneFileUpload(filenameWithoutPrefix string, content []byte) error {
	filename := s.From.FilesPath + filenameWithoutPrefix
	s.From.Files.Mut.Lock()
	// 修改Files,FilesMap
	f, ok := s.From.Files.FilesMap[filenameWithoutPrefix]
	if !ok { // 文件不存在于Files中
		s.From.Files.FilesMap[filenameWithoutPrefix] = NewFile(filenameWithoutPrefix, s.From.Files)
		f = s.From.Files.FilesMap[filenameWithoutPrefix]
	}
	f.Status = def.FILE_USING_WRITING
	f.Mut.Lock()
	defer f.Mut.Unlock()
	s.From.Files.Mut.Unlock()

	// 判断文件是否存在和是否为目录
	_, err := os.Stat(filename)
	if err != nil {
		s.From.Logger.Printf("SyncOne os.Stat-1 error, filename: %v, error: %v\n", filename, err)
	}
	targetFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		s.From.Logger.Printf("SyncOne open file error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.OPEN_FILE_ERROR)
	}
	defer targetFile.Close()
	fileContent := bytes.NewReader(content)
	n, err := io.Copy(targetFile, fileContent)
	if err != nil {
		s.From.Logger.Printf("SyncOne copy req body to file error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.WRITE_FILE_ERROR)
	}
	s.From.Logger.Printf("write %d bytes to file %s\n", n, filename)
	f.Status = def.FILE_USING_FREE
	// 更新文件info
	info, err := os.Stat(filename)
	if err != nil {
		s.From.Logger.Printf("SyncOne os.Stat-2 error, filename: %v, error: %v\n", filename, err)
	}
	f.Info = &FileInfo{
		FileName: filenameWithoutPrefix,
		Size:     info.Size(),
		CTime:    info.ModTime().UnixMilli(),
		MTime:    info.ModTime().UnixMilli(),
	}
	return nil
}

func (s *SyncServer) SyncOneBigUpload(filenameWithoutPrefix string, content []byte, first, last bool) error {
	filename := s.From.FilesPath + filenameWithoutPrefix
	s.From.Files.Mut.Lock()
	// 修改Files,FilesMap
	f, ok := s.From.Files.FilesMap[filenameWithoutPrefix]
	if !ok { // 文件不存在于Files中
		s.From.Files.FilesMap[filenameWithoutPrefix] = NewFile(filenameWithoutPrefix, s.From.Files)
		f = s.From.Files.FilesMap[filenameWithoutPrefix]
	}
	f.Status = def.FILE_USING_WRITING
	f.Mut.Lock()
	defer f.Mut.Unlock()
	s.From.Files.Mut.Unlock()

	// 判断文件是否存在和是否为目录
	_, err := os.Stat(filename)
	if err != nil {
		s.From.Logger.Printf("SyncOne os.Stat-1 error, filename: %v, error: %v\n", filename, err)
	}
	var targetFile *os.File
	if first {
		targetFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		// targetFile, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	} else {
		targetFile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, os.ModePerm)
	}
	if err != nil {
		s.From.Logger.Printf("SyncOne open file error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.OPEN_FILE_ERROR)
	}
	defer targetFile.Close()
	fileContent := bytes.NewReader(content)
	n, err := io.Copy(targetFile, fileContent)
	// targetFile.Seek(0, 2)
	// n, err := targetFile.Write(content)
	if err != nil {
		s.From.Logger.Printf("SyncOne copy req body to file error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.WRITE_FILE_ERROR)
	}
	s.From.Logger.Printf("write %d bytes to file %s\n", n, filename)
	f.Status = def.FILE_USING_FREE
	if last {
		// 更新文件info
		info, err := os.Stat(filename)
		if err != nil {
			s.From.Logger.Printf("SyncOne os.Stat-2 error, filename: %v, error: %v\n", filename, err)
		}
		f.Info = &FileInfo{
			FileName: filenameWithoutPrefix,
			Size:     info.Size(),
			CTime:    info.ModTime().UnixMilli(),
			MTime:    info.ModTime().UnixMilli(),
		}
	}
	return nil
}

func (s *SyncServer) SyncOneFileDelete(filenameWithoutPrefix string) error {
	filename := s.From.FilesPath + filenameWithoutPrefix

	s.From.Files.Mut.Lock()
	defer s.From.Files.Mut.Unlock()
	// 修改Files,FilesMap
	f, exist := s.From.Files.FilesMap[filenameWithoutPrefix]
	if !exist { // 文件不存在于Files中
		s.From.Logger.Printf("file not exist in Files, filename: %v\n", filenameWithoutPrefix)
		return errors.New(def.FILE_NOT_EXIST_IN_FILES)
	}
	f.Status = def.FILE_USING_WRITING
	f.Mut.Lock()
	defer f.Mut.Unlock()

	// 判断文件是否存在和是否为目录
	info, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		s.From.Logger.Printf("delete os.Stat error, filename: %v, error: %v\n", filename, err)
	} else if os.IsNotExist(err) {
		s.From.Logger.Printf("delete os.Stat error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.FILE_NOT_EXIST)
	}
	if info.IsDir() {
		s.From.Logger.Printf("file is a dir, filename: %v, error: %v\n", filename, err)
		return errors.New(def.FILE_IS_DIR)
	}
	err = os.Remove(filename)
	if err != nil {
		s.From.Logger.Printf("remove file error, filename: %v, error: %v\n", filename, err)
		return err
	}
	delete(s.From.Files.FilesMap, filenameWithoutPrefix)
	return nil
}

// 判断def.SYNC_RETRYING / def.SYNC_RETRY_TOO_MANY
func (c *SyncClient) SyncOneUpload(filenameWithoutPrefix string) error {
	c.From.Logger.Printf("Syncing one file, filename: %v, OP: %v\n", filenameWithoutPrefix, def.OP_UPLOAD)
	filename := c.From.FilesPath + filenameWithoutPrefix
	c.From.Files.Mut.RLock()
	// 修改Files,FilesMap
	f, exist := c.From.Files.FilesMap[filenameWithoutPrefix]
	if !exist { // 文件不存在于Files中
		c.From.Logger.Printf("file not exist in Files, filename: %v\n", filenameWithoutPrefix)
		c.From.Files.Mut.RUnlock()
		return errors.New(def.FILE_NOT_EXIST_IN_FILES)
	}
	f.Mut.RLock()
	c.From.Files.Mut.RUnlock()
	defer f.Mut.RUnlock()

	// 判断文件是否存在和是否为目录
	info, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		c.From.Logger.Printf("download os.Stat error, filename: %v, error: %v\n", filename, err)
	} else if os.IsNotExist(err) {
		c.From.Logger.Printf("download os.Stat error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.FILE_NOT_EXIST)
	}
	if info.IsDir() {
		c.From.Logger.Printf("file is a dir, filename: %v, error: %v\n", filename, err)
		return errors.New(def.FILE_IS_DIR)
	}

	f.Info = &FileInfo{
		FileName: filename,
		Size:     info.Size(),
		CTime:    info.ModTime().UnixMilli(),
		MTime:    info.ModTime().UnixMilli(),
	}

	ff, err := os.Open(filename)
	if err != nil {
		c.From.Logger.Println("SyncOne os.Open file error, ", err)
		return errors.New(def.OPEN_FILE_ERROR)
	}
	content, err := io.ReadAll(ff)
	if err != nil {
		c.From.Logger.Println("SyncOne io.ReadAll file error, ", err)
		return errors.New(def.READ_FILE_ERROR)
	}
	retry := 0
	for retry < MaxSyncOneRetry {
		ok := c.SyncOneUploadInner(filenameWithoutPrefix, content)
		if !ok {
			if retry > MaxSyncOneRetry {
				c.From.Logger.Printf("Sync one upload retry too many times, filename: %v\n", filename)
				return errors.New(def.SYNC_RETRY_TOO_MANY)
			}
			c.From.Logger.Printf("Sync one upload error, retrying, filename: %v\n", filename)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	c.From.Logger.Printf("Synced one file, filename: %v\n", filename)
	return nil
}

func (c *SyncClient) SyncOneUploadInner(filenameWithoutPrefix string, content []byte) bool {
	if !c.Connected {
		cli, err := rpc.DialHTTP("tcp", c.From.SlaveAddr.Address())
		if err != nil {
			fmt.Println(err)
			return false
		}
		c.cli = cli
		c.Connected = true
	}

	rep := SyncOneReply{false}
	args := &SyncOneArgs{
		Operation:   def.OP_UPLOAD,
		Filename:    filenameWithoutPrefix,
		NeedContent: true,
		FileContent: content,
	}
	c.cli.Call("SyncServer.SyncOne", args, &rep)
	// fmt.Println(rep.OK)
	return rep.OK
}

// 判断def.SYNC_RETRYING / def.SYNC_RETRY_TOO_MANY
func (c *SyncClient) SyncOneBigUpload(filenameWithoutPrefix string) error {
	c.From.Logger.Printf("Syncing one BIG file, filename: %v, OP: %v\n", filenameWithoutPrefix, def.OP_UPLOAD)
	filename := c.From.FilesPath + filenameWithoutPrefix
	c.From.Files.Mut.RLock()
	// 修改Files,FilesMap
	f, exist := c.From.Files.FilesMap[filenameWithoutPrefix]
	if !exist { // 文件不存在于Files中
		c.From.Logger.Printf("file not exist in Files, filename: %v\n", filenameWithoutPrefix)
		c.From.Files.Mut.RUnlock()
		return errors.New(def.FILE_NOT_EXIST_IN_FILES)
	}
	f.Mut.RLock()
	c.From.Files.Mut.RUnlock()
	defer f.Mut.RUnlock()

	// 判断文件是否存在和是否为目录
	info, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		c.From.Logger.Printf("SyncOneBig os.Stat error, filename: %v, error: %v\n", filename, err)
	} else if os.IsNotExist(err) {
		c.From.Logger.Printf("SyncOneBig os.Stat error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.FILE_NOT_EXIST)
	}
	if info.IsDir() {
		c.From.Logger.Printf("file is a dir, filename: %v, error: %v\n", filename, err)
		return errors.New(def.FILE_IS_DIR)
	}

	f.Info = &FileInfo{
		FileName: filename,
		Size:     info.Size(),
		CTime:    info.ModTime().UnixMilli(),
		MTime:    info.ModTime().UnixMilli(),
	}

	////////
	var chunknum = int64(math.Ceil(float64(info.Size()) / def.RPCSyncFileChunkSize))
	if chunknum == 1 {
		return c.SyncOneUpload(filenameWithoutPrefix)
	}
	ff, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		c.From.Logger.Println("SyncOneBig os.Open file error, ", err)
		return errors.New(def.OPEN_FILE_ERROR)
	}
	defer ff.Close()
	////////
	buffer := make([]byte, def.RPCSyncFileChunkSize)
	var i int64 = 1
	for ; i <= chunknum; i++ {
		ff.Seek((i-1)*def.RPCSyncFileChunkSize, 0)
		if len(buffer) > int(info.Size()-(i-1)*def.RPCSyncFileChunkSize) {
			buffer = make([]byte, info.Size()-(i-1)*def.RPCSyncFileChunkSize)
		}
		ff.Read(buffer)
		retry := 0
		ok, first, last := false, true, false
		if i == 1 {
			first, last = true, false
		} else if i == chunknum {
			first, last = false, true
		} else {
			first, last = false, false
		}
		for retry < MaxSyncOneRetry {
			ok = c.SyncOneBigUploadInner(filenameWithoutPrefix, buffer, chunknum, i, first, last)
			if !ok {
				if retry > MaxSyncOneRetry {
					c.From.Logger.Printf("Sync one upload retry too many times, filename: %v\n", filename)
					return errors.New(def.SYNC_RETRY_TOO_MANY)
				}
				c.From.Logger.Printf("Sync one upload error, retrying, filename: %v\n", filename)
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	c.From.Logger.Printf("Synced one BIG file, filename: %v\n", filename)
	return nil
}

func (c *SyncClient) SyncOneBigUploadInner(filenameWithoutPrefix string, content []byte, chunkNum, which int64, first, last bool) bool {
	if !c.Connected {
		cli, err := rpc.DialHTTP("tcp", ":9988")
		if err != nil {
			fmt.Println(err)
			return false
		}
		c.cli = cli
		c.Connected = true
	}

	rep := SyncOneReply{false}
	args := &SyncOneArgs{
		Operation:   def.OP_UPLOAD_BIG,
		Filename:    filenameWithoutPrefix,
		NeedContent: true,
		FileContent: content,
		IsBig:       true,
		ChunkNum:    chunkNum,
		First:       first,
		Last:        last,
		Which:       which,
	}
	c.cli.Call("SyncServer.SyncOne", args, &rep)
	// fmt.Println(rep.OK)
	return rep.OK
}

// 判断def.SYNC_RETRYING / def.SYNC_RETRY_TOO_MANY
func (c *SyncClient) SyncOneDelete(filenameWithoutPrefix string) error {
	filename := c.From.FilesPath + filenameWithoutPrefix
	// err = os.Remove(filename)
	retry := 0
	for retry < MaxSyncOneRetry {
		ok := c.SyncOneDeleteInner(filenameWithoutPrefix)
		if !ok {
			if retry > MaxSyncOneRetry {
				c.From.Logger.Printf("Sync one delete retry too many times, filename: %v\n", filename)
				return errors.New(def.SYNC_RETRY_TOO_MANY)
			}
			c.From.Logger.Printf("Sync one delete error, retrying, filename: %v\n", filename)
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	c.From.Files.Mut.Lock()
	delete(c.From.Files.FilesMap, filenameWithoutPrefix)
	c.From.Files.Mut.Unlock()
	return nil
}

func (c *SyncClient) SyncOneDeleteInner(filenameWithoutPrefix string) bool {
	if !c.Connected {
		cli, err := rpc.DialHTTP("tcp", ":9988")
		if err != nil {
			fmt.Println(err)
			return false
		}
		c.cli = cli
		c.Connected = true
	}

	rep := SyncOneReply{false}
	args := &SyncOneArgs{
		Operation:   def.OP_DELETE,
		Filename:    filenameWithoutPrefix,
		NeedContent: false,
		FileContent: []byte{},
	}
	c.cli.Call("SyncServer.SyncOne", args, &rep)
	// fmt.Println(rep.OK)
	return rep.OK
}

// 全量同步
type SyncAllArgs struct {
	Operation    string
	Filenames    []string
	NeedContent  bool
	FileContents [][]byte
}

type SyncAllReply struct {
	OKs []bool
}

func (s *SyncServer) SyncAll(args SyncAllArgs, reply *SyncAllReply) error {
	*reply = SyncAllReply{
		OKs: []bool{true, true},
	}
	fmt.Println(args.Operation, args.Filenames, args.NeedContent, string(args.FileContents[0]))
	return nil
}

func (c *SyncClient) SyncAll() bool {
	c.From.FullSyncMut.Lock()
	defer c.From.FullSyncMut.Unlock()
	flist, err := os.ReadDir(c.From.FilesPath)
	if err != nil {
		c.From.Logger.Printf("server init: os.ReadDir error: %v\n", err)
		return false
	}
	type Finfo struct {
		Filename string
		Size     int64
	}
	files := make([]Finfo, 0, len(flist))
	for _, f := range flist {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			c.From.Logger.Printf("server init: f.Info() error: %v\n", err)
			continue
		}
		files = append(files, Finfo{
			Filename: f.Name(),
			Size:     info.Size(),
		})
	}
	allGood := true
	for i := range files {
		if files[i].Size >= def.RPCSyncFileChunkSize {
			err = c.SyncOneBigUpload(files[i].Filename)
		} else {
			err = c.SyncOneUpload(files[i].Filename)
		}
		if err != nil {
			c.From.Logger.Printf("SyncAll has one error, filename: %v, %v\n", files[i].Filename, err)
			allGood = false
		}
	}
	return allGood
}

func (s *Server) RunSyncServer() {
	if IsMaster {
		// if master
		// 向slave发送sync任务，该部分代码由file调用
		// 保活检测
		s.SyncCli.SyncAll()
		s.SyncCli.RunAliveCheck()
		// return
	} else {
		// if slave
		// go 启动RPC服务，接收sync任务
		s.SyncServ.Run(s.MasterAddr.Address())
	}
}
