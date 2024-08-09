/*=============
 Author: chaidaxuan chaidaxuan@wps.cn
 Date: 2024-08-09 09:17:24
 LastEditors: chaidaxuan chaidaxuan@wps.cn
 LastEditTime: 2024-08-09 09:17:26
 FilePath: /Golang-Samples/src/_Utils/Log/logger.go
 Description:



 Copyright (c) 2024 by chaidaxuan, All Rights Reserved.
=============*/

package log_test

import (
	"fmt"
	"log"
	"os"
)

type ULog struct {
	Path  string
	l     *log.Logger
	print bool
}

const LOG_PRINT_TO_CONSOLE = true

func NewLogger(path string) *ULog {
	logfile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &ULog{
		Path:  path,
		l:     log.New(logfile, "", log.LstdFlags),
		print: LOG_PRINT_TO_CONSOLE,
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
