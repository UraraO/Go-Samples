#! /bin/sh
###
 # @Author: chaidaxuan chaidaxuan@wps.cn
 # @Date: 2024-07-31 17:54:42
 # @LastEditors: chaidaxuan chaidaxuan@wps.cn
 # @LastEditTime: 2024-08-01 17:23:20
 # @FilePath: /urarao/GoProjects/chaidaxuan/filesync/build.sh
 # @Description: 
 # 
 # Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
### 
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sync
