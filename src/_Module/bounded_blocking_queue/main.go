/*===========
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:47:01
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-06 22:42:12
 FilePath: /Golang-Samples/src/_Module/bounded_blocking_queue/main.go
 Description:

 该库为一个并发工具使用例程库，concurrency文件夹中有更多示例；
 cache是一个并发安全的K-V键值对缓存工具
 singleton是一个单例模式的golang实现
 ccrctest是一个全局唯一ID生成器，实际运行中可以使用uuid库
 spinlock是一个自旋锁的golang实现
 bank是一个并发安全的银行模块实现，可以多账号转账，存取款等

 Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
===========*/

package main

import "Golang-Samples/src/_Module/bounded_blocking_queue/concurrency"

func main() {
	concurrency.BlockQueueTest()
}
