#! /bin/sh
###
 # @Author: chaidaxuan chaidaxuan@wps.cn
 # @Date: 2024-07-31 17:54:36
 # @LastEditors: chaidaxuan chaidaxuan@wps.cn
 # @LastEditTime: 2024-08-01 11:07:42
 # @FilePath: /urarao/GoProjects/chaidaxuan/filesync/master.sh
 # @Description: 
 # 
 # Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
### 
./sync master ./files_master --http-listen-addr=127.0.0.1 --http-listen-port=9988 --slave-grpc-addr=127.0.0.1 --slave-grpc-port=9989 --server-side-encryption false
