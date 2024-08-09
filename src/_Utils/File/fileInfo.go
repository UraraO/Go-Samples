/*=============
 Author: chaidaxuan chaidaxuan@wps.cn
 Date: 2024-08-09 10:35:40
 LastEditors: chaidaxuan chaidaxuan@wps.cn
 LastEditTime: 2024-08-09 10:50:12
 FilePath: /Golang-Samples/src/_Utils/File/fileInfo.go
 Description:



 Copyright (c) 2024 by chaidaxuan, All Rights Reserved.
=============*/

package file_utils

import (
	"io/fs"
	"os"
	"runtime"
	"syscall"
	"time"
)

// 获取文件创建时间
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
// 该函数仅可在linux环境下使用
func GetFileCreateTime(filePath string) time.Time {
	// 获取文件原来的访问时间，修改时间
	finfo, err := os.Stat(filePath)
	if err != nil {
		return time.Unix(0, 0)
	}

	linuxFileAttr := finfo.Sys().(*syscall.Stat_t)

	if runtime.GOOS == "linux" {
		return SecondToTime(linuxFileAttr.Ctim.Sec)
	} else if runtime.GOOS == "windows" {
		return time.Unix(0, 0)
	}
	return time.Now()
}

// 判断文件是否存在
func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

// 判断文件是否存在和是否为目录
func CheckFileExistorIsDir(filepath string) (exist bool, isdir bool, info fs.FileInfo) {
	info, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) { // 文件不存在
			return false, false, nil
		} else { // 未知错误（type *PathError）
			return false, false, nil
		}
	}
	if info.IsDir() {
		return true, true, info
	}
	return true, false, info
}
