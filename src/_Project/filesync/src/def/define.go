/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-29 14:09:56
  - @LastEditors: chaidaxuan chaidaxuan@wps.cn
  - @LastEditTime: 2024-08-01 17:15:20
  - @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/def/define.go

* @Description:

全局预定义常量，包含各类错误，类型状态，操作类型等
实际运行中，可以通过配置文件，环境变量等传递

* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package def

const RPCSyncFileChunkSize = 32 * (1 << 10)
const DEFAULT_LOG_PATH = "./logs/"
const DEFAULT_LOG_NAME_MASTER = "server-master.log"
const DEFAULT_LOG_NAME_SLAVE = "server-slave.log"
const LOG_PRINT_TO_CONSOLE = false
const MIDFILE_SUFFIX = ".midxx"

// Server Role
const (
	ROLE_INIT = iota
	ROLE_MASTER
	ROLE_SLAVE
)

// File Status
const (
	FILE_USING_INIT = iota
	FILE_USING_FREE
	FILE_USING_BUSY
	FILE_USING_CREATING
	FILE_USING_READING
	FILE_USING_WRITING
)

// Mission Type
const (
	MISSION_TYPE_INIT = iota
	MISSION_TYPE_READ
	MISSION_TYPE_WRITE
	MISSION_TYPE_UPLOAD
	MISSION_TYPE_DOWNLOAD
	MISSION_TYPE_LIST
	MISSION_TYPE_DELETE
)

// Sync operation
const (
	OP_UPLOAD     = "upload"
	OP_DELETE     = "delete"
	OP_UPLOAD_BIG = "upload_big"
)

// ERROR
const (
	OPEN_FILE_ERROR         = "OPEN_FILE_ERROR"
	WRITE_FILE_ERROR        = "WRITE_FILE_ERROR"
	DELETE_FILE_ERROR       = "DELETE_FILE_ERROR"
	RENAME_FILE_ERROR       = "RENAME_FILE_ERROR"
	READ_FILE_ERROR         = "READ_FILE_ERROR"
	FILE_IS_DIR             = "FILE_IS_DIR"
	FILE_NOT_EXIST          = "FILE_NOT_EXIST"
	FILE_NOT_EXIST_IN_FILES = "FILE_NOT_EXIST_IN_FILES"

	SYNC_RETRY_TOO_MANY = "SYNC_RETRY_TOO_MANY"
	SYNC_RETRYING       = "SYNC_RETRYING"
)
