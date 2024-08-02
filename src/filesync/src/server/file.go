/*
  - @Author: chaidaxuan chaidaxuan@wps.cn
  - @Date: 2024-07-29 15:13:17
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-08-02 10:54:12
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/server/file.go
  - @Description:

文件服务器的文件服务部分，提供了文件上传（用户上传文件到服务器），文件下载（从服务器下载文件），文件列表（获取服务器文件列表）和文件删除（删除服务器文件）共4个接口，由http服务调用

  - Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package server

import (
	"encoding/json"
	"errors"
	"filesync/src/def"
	"filesync/src/utils"
	"io"
	"net/http"
	"os"
	"sync"
	"syscall"
)

type File struct {
	FileName string
	Status   int
	Info     *FileInfo
	// 删除操作结束后，使用From从list中移除 *File 自身
	From *FileList

	// 该互斥量用于保护单个文件
	// 加锁和解锁操作用于读取文件/修改文件
	// 读锁：列表查询/下载；加锁，下载（向客户端发送），解锁
	// 写锁：上传/删除；加锁，文件上传和删除，解锁
	Mut sync.RWMutex
}

type FileInfo struct {
	FileName string `json:"name"`
	Size     int64  `json:"size"`
	CTime    int64  `json:"c_time"`
	MTime    int64  `json:"m_time"`
}

type Files struct { // 仅用于List API生成JSON
	Files []*FileInfo `json:"files"`
}

type FileList struct {
	FilesMap map[string]*File

	// 该互斥量用于保护整个文件列表
	// 加锁和解锁操作仅用于修改文件列表，与文件本身的读写操作无关
	// 读锁：列表查询/下载，获取到 *File 文件信息后立刻解锁
	// 写锁：上传/删除，文件上传和删除操作完成后，加锁，修改，立刻解锁
	Mut sync.RWMutex
}

// File part
func NewFile(fileName string, list *FileList) *File {
	return &File{
		FileName: fileName,
		Status:   def.FILE_USING_INIT,
		From:     list,
		Mut:      sync.RWMutex{},
	}
}

// FileList part
func NewFileList() *FileList {
	return &FileList{
		FilesMap: make(map[string]*File, 10),
		Mut:      sync.RWMutex{},
	}
}

// file server
func (s *Server) RunFileServer() {
	// go 文件服务，消费任务，从s.MissionQueue取出任务

	// master / slave
	// read: List / Download

	// master
	// write: Upload / Delete

}

// Upload part
func (s *Server) UploadFile(filenameWithoutPrefix string, fileContent io.ReadCloser) error {
	s.FullSyncMut.RLock()
	defer s.FullSyncMut.RUnlock()
	filename := s.FilesPath + filenameWithoutPrefix
	s.Files.Mut.Lock()
	// 修改Files,FilesMap
	f, ok := s.Files.FilesMap[filenameWithoutPrefix]
	if !ok { // 文件不存在于Files中
		s.Files.FilesMap[filenameWithoutPrefix] = NewFile(filenameWithoutPrefix, s.Files)
		f = s.Files.FilesMap[filenameWithoutPrefix]
	}
	f.Status = def.FILE_USING_WRITING
	f.Mut.Lock()
	defer f.Mut.Unlock()
	s.Files.Mut.Unlock()

	// 判断文件是否存在和是否为目录
	// finfo, err := os.Stat(filename)
	// exist := false
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		s.Logger.Printf("upload os.Stat-1 error, filename: %v, error: %v\n", filename, err)
	// 	} else if finfo.IsDir() {
	// 		s.Logger.Printf("upload file is a dir, filename: %v, error: %v\n", filename, err)
	// 		return errors.New(def.FILE_IS_DIR)
	// 	}
	// } else {
	// 	exist = true
	// }
	exist, isdir, info := utils.CheckFileExistorIsDir(filename)
	if isdir {
		s.Logger.Printf("upload file is a dir, filename: %v, error: %v\n", filename, errors.New(def.FILE_IS_DIR))
		return errors.New(def.FILE_IS_DIR)
	}

	midFilename := filename + def.MIDFILE_SUFFIX
	midFile, err := os.OpenFile(midFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	// targetFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		s.Logger.Printf("upload open file error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.OPEN_FILE_ERROR)
	}
	defer midFile.Close()
	// defer targetFile.Close()
	n, err := io.Copy(midFile, fileContent)
	if err != nil {
		s.Logger.Printf("upload copy req body to file error, filename: %v, error: %v\n", filename, err)
		os.Remove(midFilename)
		return errors.New(def.WRITE_FILE_ERROR)
	}
	if exist {
		err = os.Remove(filename)
		if err != nil {
			s.Logger.Printf("upload os.Remove error, filename: %v, error: %v\n", filename, err)
			return errors.New(def.DELETE_FILE_ERROR)
		}
	}
	err = os.Rename(midFilename, filename)
	if err != nil {
		s.Logger.Printf("upload os.Rename error, filename: %v, error: %v\n", filename, err)
		return errors.New(def.RENAME_FILE_ERROR)
	}
	s.Logger.Printf("write %d bytes to file %s\n", n, filename)
	f.Status = def.FILE_USING_FREE
	// Sync
	go func() {
		s.FullSyncMut.RLock()
		defer s.FullSyncMut.RUnlock()
		var err error
		if n >= def.RPCSyncFileChunkSize {
			err = s.SyncCli.SyncOneBigUpload(filenameWithoutPrefix)
		} else {
			err = s.SyncCli.SyncOneUpload(filenameWithoutPrefix)
		}
		if err != nil {
			if err.Error() == def.SYNC_RETRYING {
				s.Logger.Printf("SyncOne Upload retrying, file: %v\n", filenameWithoutPrefix)
			} else if err.Error() == def.SYNC_RETRY_TOO_MANY {
				s.Logger.Printf("SyncOne Upload retryed TOO many times, file: %v\n", filenameWithoutPrefix)
			}
		}
	}()
	// 更新文件info
	info, err = os.Stat(filename)
	if err != nil {
		s.Logger.Printf("upload os.Stat-2 error, filename: %v, error: %v\n", filename, err)
	}
	f.Info = &FileInfo{
		FileName: filenameWithoutPrefix,
		Size:     info.Size(),
		CTime:    info.ModTime().UnixMilli(),
		MTime:    info.ModTime().UnixMilli(),
	}
	return nil
}

// Download part
func (s *Server) DownloadFile(w http.ResponseWriter, r *http.Request, filenameWithoutPrefix string) error {
	s.FullSyncMut.RLock()
	defer s.FullSyncMut.RUnlock()
	filename := s.FilesPath + filenameWithoutPrefix
	s.Files.Mut.RLock()
	// 修改Files,FilesMap
	f, exist := s.Files.FilesMap[filenameWithoutPrefix]
	if !exist { // 文件不存在于Files中
		s.Logger.Printf("file not exist in Files, filename: %v\n", filenameWithoutPrefix)
		w.WriteHeader(http.StatusNotFound)
		s.Files.Mut.RUnlock()
		return errors.New(def.FILE_NOT_EXIST_IN_FILES)
	}
	f.Mut.RLock()
	s.Files.Mut.RUnlock()
	defer f.Mut.RUnlock()

	// 判断文件是否存在和是否为目录
	// info, err := os.Stat(filename)
	// if err != nil {
	// 	// s.Logger.Printf("download os.Stat error, filename: %v, error: %v\n", filename, err)
	// 	s.Logger.Printf("download os.Stat error, filename: %v, error: %v\n", filename, err)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return errors.New(def.FILE_NOT_EXIST)
	// }
	// if info.IsDir() {
	// 	s.Logger.Printf("file is a dir, filename: %v, error: %v\n", filename, err)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return errors.New(def.FILE_IS_DIR)
	// }
	exist, isdir, info := utils.CheckFileExistorIsDir(filename)
	if isdir {
		s.Logger.Printf("file is a dir, filename: %v\n", filename)
		w.WriteHeader(http.StatusNotFound)
		return errors.New(def.FILE_IS_DIR)
	}
	if !exist {
		s.Logger.Printf("download file not exist, filename: %v\n", filename)
		w.WriteHeader(http.StatusNotFound)
		return errors.New(def.FILE_NOT_EXIST)
	}
	// info, err := os.Stat(filename)
	// if err != nil {
	// 	s.Logger.Printf("download os.Stat error, filename: %v, error: %v\n", filename, err)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return errors.New(def.FILE_NOT_EXIST)
	// }
	// f.Info = &FileInfo{
	// 	FileName: filename,
	// 	Size:     info.Size(),
	// 	CTime:    info.ModTime().UnixMilli(),
	// 	MTime:    info.ModTime().UnixMilli(),
	// }

	w.Header().Set("Last-Modified", info.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	http.ServeFile(w, r, filename)
	return nil
}

// Delete part
func (s *Server) DeleteFile(filenameWithoutPrefix string) error {
	s.FullSyncMut.RLock()
	defer s.FullSyncMut.RUnlock()
	filename := s.FilesPath + filenameWithoutPrefix
	s.Files.Mut.Lock()
	defer s.Files.Mut.Unlock()
	// 修改Files,FilesMap
	f, exist := s.Files.FilesMap[filenameWithoutPrefix]
	if !exist { // 文件不存在于Files中
		s.Logger.Printf("file not exist in Files, filename: %v\n", filenameWithoutPrefix)
		return errors.New(def.FILE_NOT_EXIST_IN_FILES)
	}
	f.Status = def.FILE_USING_WRITING
	f.Mut.Lock()
	defer f.Mut.Unlock()

	// 判断文件是否存在和是否为目录
	// info, err := os.Stat(filename)
	// if err != nil {
	// 	// s.Logger.Printf("delete os.Stat error, filename: %v, error: %v\n", filename, err)
	// 	s.Logger.Printf("delete os.Stat error, filename: %v, error: %v\n", filename, err)
	// 	return errors.New(def.FILE_NOT_EXIST)
	// }
	// if info.IsDir() {
	// 	s.Logger.Printf("file is a dir, filename: %v, error: %v\n", filename, err)
	// 	return errors.New(def.FILE_IS_DIR)
	// }
	exist, isdir, _ := utils.CheckFileExistorIsDir(filename)
	if isdir {
		s.Logger.Printf("file is a dir, filename: %v\n", filename)
		return errors.New(def.FILE_IS_DIR)
	}
	if !exist {
		s.Logger.Printf("download file not exist, filename: %v\n", filename)
		return errors.New(def.FILE_NOT_EXIST)
	}
	err := os.Remove(filename)
	if err != nil {
		s.Logger.Printf("remove file error, filename: %v, error: %v\n", filename, err)
		return err
	}
	delete(s.Files.FilesMap, filenameWithoutPrefix)
	// Sync
	go func() {
		s.FullSyncMut.RLock()
		defer s.FullSyncMut.RUnlock()
		err := s.SyncCli.SyncOneDelete(filenameWithoutPrefix)
		if err != nil {
			if err.Error() == def.SYNC_RETRYING {
				s.Logger.Printf("SyncOne Upload retrying, file: %v\n", filenameWithoutPrefix)
			} else if err.Error() == def.SYNC_RETRY_TOO_MANY {
				s.Logger.Printf("SyncOne Upload retryed TOO many times, file: %v\n", filenameWithoutPrefix)
			}
		}
	}()
	return nil
}

// List part
func (s *Server) GetFilelistJSON(limit, offset int) []byte {
	files := s.GetFilelist(limit, offset)
	b, err := json.Marshal(files)
	if err != nil {
		s.Logger.Println("GetFilelistJSON error:", err)
		return []byte{}
	}
	return b
}

func (s *Server) GetFilelist(limit, offset int) *Files {
	s.FullSyncMut.RLock()
	defer s.FullSyncMut.RUnlock()
	s.Files.Mut.RLock()
	defer s.Files.Mut.RUnlock()
	size := len(s.Files.FilesMap)
	flist := Files{
		Files: make([]*FileInfo, 0, size),
	}
	// 读取全部文件列表
	if offset > size || limit == 0 {
		return &flist
	}
	i, num := 0, 0
	for fname, f := range s.Files.FilesMap {
		// 获取文件信息
		if i < offset {
			i++
			continue
		}
		if num >= limit {
			break
		}
		f.Mut.RLock()
		info, err := s.getFileInfo(s.FilesPath + fname)
		if err != nil {
			s.Logger.Printf("getFileInfo error, filename: %v, error: %v\n", fname, err)
			f.Mut.RUnlock()
			continue
		}
		info.FileName = fname
		flist.Files = append(flist.Files, info)
		num++
		f.Mut.RUnlock()
	}
	return &flist
}

func (s *Server) getFileInfo(fname string) (res *FileInfo, err error) {
	info, err := os.Stat(fname)
	if err != nil {
		s.Logger.Printf("getFileInfo error, filename: %v, error: %v\n", fname, err)
		return nil, err
	}
	if info.IsDir() {
		return nil, errors.New(def.FILE_IS_DIR)
	}
	linuxFileAttr := info.Sys().(*syscall.Stat_t)
	res = &FileInfo{
		FileName: fname,
		Size:     linuxFileAttr.Size,
		CTime:    linuxFileAttr.Ctim.Sec,
		MTime:    linuxFileAttr.Mtim.Sec,
	}
	return res, nil
}
