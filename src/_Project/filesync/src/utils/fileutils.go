/*
  - @Author: chaidaxuan chaidaxuan@wps.cn
  - @Date: 2024-07-29 16:42:23
  - @LastEditors: chaidaxuan chaidaxuan@wps.cn
  - @LastEditTime: 2024-08-02 10:34:46
  - @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/utils/fileutils.go
  - @Description:

定义了一些文件操作的辅助方法

  - Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"syscall"
	"time"
)

// 秒级时间戳转为时间类型
func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func MilliSecondToTime(msec int64) time.Time {
	return time.UnixMilli(msec)
}

func MicroSecondToTime(msec int64) time.Time {
	return time.UnixMicro(msec)
}

// 获取文件创建时间
func GetFileCreateTime(filePath string) time.Time {
	// 获取文件原来的访问时间，修改时间
	finfo, _ := os.Stat(filePath)

	// linux环境下代码如下
	linuxFileAttr := finfo.Sys().(*syscall.Stat_t)

	if runtime.GOOS == "linux" {
		return SecondToTime(linuxFileAttr.Ctim.Sec)
	} else if runtime.GOOS == "windows" {
		return time.Now()
	}
	return time.Now()
}

// 判断文件是否存在和是否为目录
func CheckFileExistorIsDir(filepath string) (exist bool, isdir bool, info fs.FileInfo) {
	info, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		} else {
			return false, false, nil
		}
	} else {
		if info.IsDir() {
			return true, true, info
		}
		return true, false, info
	}
}

func CreateDir(dirPath string) error {
	// 判断文件是否存在和是否为目录
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dirPath, os.ModePerm)
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return errors.New("there is a file has same name")
	}
	return nil
}

func GenerateBigFile(path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i <= 1000000; i++ {
		_, err = f.Write([]byte(fmt.Sprintf("%d ", i)))
		if err != nil {
			panic(err)
		}
	}
}

// func GetFilelistJSON(flist []string) []byte {
// 	files := GetFilelist(flist)
// 	b, err := json.Marshal(files)
// 	if err != nil {
// 		fmt.Println("GetFilelistJSON error:", err)
// 		return []byte{}
// 	}
// 	fmt.Println(string(b))
// 	return b
// }

// func GetFilelist(FilesMap []string) *Files {
// 	flist := Files{
// 		Files: make([]*FileInfo, 0, 10),
// 	}

// 	// 读取全部文件列表
// 	for _, fname := range FilesMap {
// 		// 获取文件信息
// 		info, err := getFileInfo(fname)
// 		if err != nil {
// 			continue
// 		}
// 		flist.Files = append(flist.Files, info)

// 	}
// 	return &flist
// }

// func getFileInfo(fname string) (res *FileInfo, err error) {
// 	info, err := os.Stat("./src/" + fname)
// 	if err != nil {
// 		fmt.Printf("getFileInfo error, filename: %v, error: %v\n", fname, err)
// 		return nil, err
// 	}
// 	linuxFileAttr := info.Sys().(*syscall.Stat_t)
// 	res = &FileInfo{
// 		FileName: fname,
// 		Size:     linuxFileAttr.Size,
// 		CTime:    linuxFileAttr.Ctim.Sec,
// 		MTime:    linuxFileAttr.Mtim.Sec,
// 	}
// 	fmt.Println("getFileInfo: ", *res)
// 	return res, nil
// }
