/*=============
 Author: chaidaxuan chaidaxuan@wps.cn
 Date: 2024-07-26 18:06:52
 LastEditors: chaidaxuan chaidaxuan@wps.cn
 LastEditTime: 2024-08-09 14:25:00
 FilePath: /Golang-Samples/src/_Utils/File/dir.go
 Description:



 Copyright (c) 2024 by chaidaxuan, All Rights Reserved.
=============*/

package file_utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	FileName string `json:"name"`
	Size     int64  `json:"size"`
	CTime    int64  `json:"c_time"`
	MTime    int64  `json:"m_time"`
}

// Mkdir 创建目录
// 假设FileServer/bin/main.exe运行Mkdir("./FileData", username)
// 则其创建目录在控制台当前所在路径的FileData文件夹下
// 即最终创建：FileServer/FileData/username-2023-7-6
// FileDataMkdir("./FileData", username)
func MkdirXXX(basePath string, username string) string {
	//	1.获取当前时间,并且格式化时间
	folderName := username + "-" + time.Now().Format("2006-01-02")
	folderPath := filepath.Join(basePath, folderName)
	// 使用mkdirall会创建多层级目录
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		fmt.Println("Mkdir MkdirAll ERROR:", err)
		return ""
	}
	return folderPath
}

func MkdirTest() {
	username := "Urara"
	MkdirXXX("FileData", username)
	os.Mkdir("./FileData/dir_test", os.ModeDir)
}

func CreateDir(dirPath string) error {
	// 判断文件是否存在和是否为目录
	exist, isdir, _ := CheckFileExistorIsDir(dirPath)
	if exist {
		if isdir {
			return fmt.Errorf("there is a dir has same name")
		}
		return fmt.Errorf("there is a file has same name")
	}
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// LoadFilesInDir 加载文件夹下的全部文件，不递归遍历子文件夹
// LoadFilesInDir("./FileData")
func LoadFilesInDir(dirPath string) ([]FileInfo, error) {
	flist, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, 0, len(flist))
	for i := range flist {
		if flist[i].IsDir() {
			continue
		}
		info, err := flist[i].Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			FileName: flist[i].Name(),
			Size:     info.Size(),
			CTime:    info.ModTime().UnixMilli(),
			MTime:    info.ModTime().UnixMilli(),
		})
	}
	return files, nil
}

// LoadFilesInDirRecursive 加载文件夹下的全部文件，递归遍历子文件夹
func LoadFilesInDirRecursive(dirPath string) ([]FileInfo, error) {
	flist, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, 0, len(flist))
	for i := range flist {
		if flist[i].IsDir() {
			fs, err := LoadFilesInDirRecursive(dirPath + "/" + flist[i].Name())
			if err != nil {
				continue
			}
			files = append(files, fs...)
			continue
		}
		info, err := flist[i].Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			FileName: flist[i].Name(),
			Size:     info.Size(),
			CTime:    info.ModTime().UnixMilli(),
			MTime:    info.ModTime().UnixMilli(),
		})
	}
	return files, nil
}
