/*
  - @Author: chaidaxuan chaidaxuan@wps.cn
  - @Date: 2024-07-29 14:56:55
  - @LastEditors: chaidaxuan chaidaxuan@wps.cn
  - @LastEditTime: 2024-08-01 14:19:27
  - @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/logger/logger.go
  - @Description:

定义了一个日志工具类，用于打印日志和写入日志文件

  - Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package logger

import (
	"filesync/src/def"
	"fmt"
	"log"
	"os"
)

type ULog struct {
	Path  string
	l     *log.Logger
	print bool
}

func NewLogger(path string) *ULog {
	logfile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &ULog{
		Path:  path,
		l:     log.New(logfile, "", log.LstdFlags),
		print: def.LOG_PRINT_TO_CONSOLE,
	}
}

func (l *ULog) Printf(format string, args ...any) {
	l.l.Printf(format, args...)
	if l.print {
		fmt.Printf(format, args...)
	}
}

func (l *ULog) Println(args ...any) {
	l.l.Println(args...)
	if l.print {
		fmt.Println(args...)
	}
}

func (l *ULog) Fatalf(format string, args ...any) {
	l.l.Fatalf(format, args...)
	if l.print {
		fmt.Printf(format, args...)
	}
}

func (l *ULog) Fatalln(args ...any) {
	l.l.Fatalln(args...)
	if l.print {
		fmt.Println(args...)
	}
}

func (l *ULog) Panicf(format string, args ...any) {
	l.l.Panicf(format, args...)
	if l.print {
		fmt.Printf(format, args...)
	}
}

func (l *ULog) Panicln(args ...any) {
	l.l.Panicln(args...)
	if l.print {
		fmt.Println(args...)
	}
}
