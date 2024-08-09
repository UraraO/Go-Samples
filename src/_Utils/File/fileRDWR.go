/*=============
 Author: chaidaxuan chaidaxuan@wps.cn
 Date: 2024-07-26 18:06:52
 LastEditors: chaidaxuan chaidaxuan@wps.cn
 LastEditTime: 2024-08-09 15:31:50
 FilePath: /Golang-Samples/src/_Utils/File/fileRDWR.go
 Description:



 Copyright (c) 2024 by chaidaxuan, All Rights Reserved.
=============*/

package file_utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func ReadFile() []byte {
	f, err := os.OpenFile("file.txt", os.O_RDWR, os.ModeTemporary)
	if err != nil {
		fmt.Println("read file fail", err)
		return []byte{}
	}
	defer f.Close()

	fd, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return []byte{}
	}

	fmt.Println(string(fd))
	return fd
}

func WriteF(src []byte, fileName string) {
	err := os.WriteFile(fileName, src, 0666)
	if err != nil {
		fmt.Println("write fail")
	}
	fmt.Println("write success")
}

// 流式读写内容，io.Copy可以直接将流拷贝到另一个读写流，无需担心内存溢出
// 该接口常用于http上传文件等
func IOWriteFileInner(src io.Reader, dst io.Writer) error {
	// midFilename := filename + def.MIDFILE_SUFFIX
	// midFile, err := os.OpenFile(midFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	// targetFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	// if err != nil {
	// 	s.Logger.Printf("upload open file error, filename: %v, error: %v\n", filename, err)
	// 	return errors.New(def.OPEN_FILE_ERROR)
	// }
	// defer midFile.Close()
	// defer targetFile.Close()
	n, err := io.Copy(dst, src)
	if err != nil {
		return errors.New("io copy error")
	}
	log.Println("write file size: ", n)
	return nil
}

// 流式读写文件
func IOWriteFile(srcFile string, dstFile string) error {
	srcF, err := os.OpenFile(srcFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer srcF.Close()
	dstF, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dstF.Close()
	// err = IOWriteFileInner(srcF, dstF)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	n, err := io.Copy(dstF, srcF)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("write file size: ", n)
	return nil
}

func FileRDWRtest() {
	fmt.Println(CheckFileExist("file.txt"))
	WriteF(ReadFile(), "file2.txt")
}

// 分块读文件
const chunkSize = 1024 * 1024 // 1MB

func ReadFileBlock(filepath string) {
	exist, isdir, fileInfo := CheckFileExistorIsDir(filepath)
	if !exist {
		fmt.Println("文件不存在")
		return
	}
	if isdir {
		fmt.Println("filepath is dir")
		return
	}
	// blockNum是分块数量
	blockNum := int64(math.Ceil(float64(fileInfo.Size()) / chunkSize))

	src, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()
	buf := make([]byte, chunkSize)
	var i int64 = 1 // 第i块,从1开始
	for ; i <= blockNum; i++ {
		src.Seek((i-1)*chunkSize, 0)
		if i == blockNum { // 最后一块
			if chunkSize > int(fileInfo.Size()-(i-1)*chunkSize) {
				buf = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
			}
		}
		n, err := src.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("block ", i, "len:", n, ":", string(buf))
	}
}

// 生成大文件
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
