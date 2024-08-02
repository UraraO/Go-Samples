#! /bin/sh
###
 # @Author: chaidaxuan chaidaxuan@wps.cn
 # @Date: 2024-07-31 17:54:42
 # @LastEditors: chaidaxuan chaidaxuan@wps.cn
 # @LastEditTime: 2024-08-01 11:07:49
 # @FilePath: /urarao/GoProjects/chaidaxuan/filesync/slave.sh
 # @Description: 
 # 
 # Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
### 
./sync slave ./files_slave --http-listen-addr=127.0.0.1 --http-listen-port=9990 --grpc-listen-addr=127.0.0.1 --grpc-listen-port=9989 --server-side-encryption false
