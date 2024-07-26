package cron

import (
	"fmt"
	"sync"
)

// cron工具帮助您快速创建定时任务
// 开始之前，请确保您已经读过该目录下的图片：Cron表达式.png
// 或参考该网址：https://zhuanlan.zhihu.com/p/437328366
// 或查看该文件最下方的内容
// //////////////////////////////////////////////////////////
// 如果您需要快速编写定时任务并执行，只需要选择下方模板并拷贝
// 您要做的：
//
//	1.更改任务名：__YourTaskName__
//	2.在Run()中实现您的定时任务
//	3.创建一个Cron：NewBasicCron() 或 NewDistributeCron(&Mutex)
//	4.选择一种Add方式，将任务添加到Cron实例中
//	5.启动Cron：c.Start()

// 模板一：基础定时任务，绝大多数情况使用该模板即可
type __YourTaskName__ struct {
	// TODO anything you want
	// TODO parameters or status
}

func (task *__YourTaskName__) Run() {
	// TODO 定时任务的业务逻辑
	// Mutex.lock() 若访问共享数据，您负责管理锁，可以在Task中保存锁
	// 如果您希望在分布式系统中使用，请考虑在此处使用分布式锁
	fmt.Println("Anything you want to do")
	// Mutex.Unlock()
}

/*
// 若您需要回调方法，则取消此处注释，并实现回调函数的Callback()方法
type __CallbackFunc__ struct {
	// anything you want
}

func (callback *CallbackFunc) Callback() {
	fmt.Println("Callback")
}

func __main__() {
	task1 := __YourTaskName__{}
	cb := callbackFunc{}
	c.AddTaskWithCallback("@every 1m", &task1, "Task 3s_1", &cb)
	// task1执行成功后，执行您的回调函数
}
*/

func __main__() {
	// 创建Task
	task1 := __YourTaskName__{}
	task2 := __YourTaskName__{}
	// 创建Cron
	c := NewBasicCron()
	// 将任务交由Cron管理
	// 添加一个简单函数，每秒执行一次
	c.AddFunc("@every 1s", func() {
		fmt.Println("task do every 1 second")
	}, "what the func will do?")
	// 添加您的自定义任务
	// 每1小时30分45秒执行一次
	c.AddTask("@every 1h30m45s", &task1, "what your task will do?")
	// 添加您的自定义任务
	// 标准Cron表达式，自定义任意间隔或定时，此处为每天凌晨4点执行
	// 两个额外option：delay = true，若当前任务未结束就到下个任务，则下个任务延迟到本次任务完成再执行（默认跳过下个任务，delay = false，而非延迟）
	// coverPanic = false，若您的Task抛出panic，则Cron不会帮您捕获异常，而是导致程序异常结束（默认捕获异常，coverPanic = true）
	c.AddTaskWithOption("0 0 4 * * ?", &task2, "what your task will do?", true, false)

	// 运行
	c.Start()
	// time.Sleep(10 * time.Second)
	// c.Stop()

	// c.ListExistTask()    列出当前Cron管理的Task信息
	// Task: {taskID  Description  Skip?  CoverPanic? }
	// c.Remove(taskID)         上方列出的任务中保存ID信息，根据ID移除任务
}

// ///////////////////////////////////////

// 模板二：分布式单实例定时任务
// 该模板需要您传入一个分布式锁的地址
// Cron会确保分布式系统中，同时只有一个Cron实例持有该锁
type __YourDistTaskName__ struct {
	// TODO anything you want
	// TODO parameters or status
}

func (task *__YourDistTaskName__) Run() {
	// TODO 定时任务的业务逻辑
	// Mutex.lock() 若访问共享数据，您负责管理锁，可以在Task中保存锁
	// 分布式的Cron实例保证同时仅有一个实例正在Start()，基于您传入的分布式锁
	fmt.Println("Anything you want to do")
	// Mutex.Unlock()
}

/*
// 若您需要回调方法，则取消此处注释，并实现回调函数的Callback()方法
type __CallbackFunc__ struct {
	// anything you want
}

func (callback *CallbackFunc) Callback() {
	fmt.Println("Callback")
}

func __distribute_main__() {
	task1 := __YourTaskName__{}
	cb := callbackFunc{}
	c.AddTaskWithCallback("@every 1m", &task1, "Task 3s_1", &cb)
	// task1执行成功后，执行您的回调函数
}
*/

func __distribute_main__() {
	// 创建Task
	task1 := __YourTaskName__{}
	task2 := __YourTaskName__{}
	// 创建Cron
	// TODO 务必将此处mutex替换为您的分布式锁，此处仅用于展示
	var TODOReplaceThisMutex sync.Mutex
	c := NewDistributeCron(&TODOReplaceThisMutex)
	// 将任务交由Cron管理
	// 添加一个简单的函数，每秒执行一次
	c.AddFunc("@every 1s", func() {
		fmt.Println("task do every 1 second")
	}, "what the task will do?")
	// 添加您的自定义任务
	// 每1小时30分45秒执行一次
	c.AddTask("@every 1h30m45s", &task1, "what your task will do?")
	// 添加您的自定义任务
	// 标准Cron表达式，自定义任意间隔或定时，此处为每天凌晨4点执行
	// 两个额外option：delay = true，若当前任务未结束就到下个任务，则下个任务延迟到本次任务完成再执行（默认跳过下个任务，delay = false，而非延迟）
	// coverPanic = false，若您的Task抛出panic，则Cron不会帮您捕获异常，而是导致程序异常结束（默认捕获异常，coverPanic = true）
	c.AddTaskWithOption("0 0 4 * * ?", &task2, "what your task will do?", true, false)

	// 运行
	c.Start()
	// time.Sleep(10 * time.Second)
	// c.Stop()

	// c.ListExistTask()    列出当前Cron管理的Task
	// Task: {taskID  Description  Skip?  CoverPanic? }
	// c.Remove(ID)         上方列出的任务中保存ID信息，根据ID移除任务
}

// ///////////////////////////////////////
// 关于Cron表达式：
/*
cron库中的cron表达式和预定义时间写法
    常用cron表达式，以下示例默认带有秒级，即：秒  分  时  日  月  周
        " * * * * * ? " :  每秒执行一次
        " 0 * * * * ? " :  每分钟执行一次（每分钟第 0 秒）
        " 0 * 9-18 * * ? " :  每天的 9 到 18 点中，每秒执行一次，限制时间在工作时间中
        " 0 0 0/6 * * ? " :  每天 0 点开始，每隔 6 小时执行一次，低频操作，如批量同步数据
        " 0 10 4 * * 0,6 " :  每周的周六周日，凌晨 4 点 10 分执行一次

    常用Cron库定义的方便写法
        预定义时间规则
            `@yearly`：也可以写作@annually，表示每年第一天的 0 点。等价于 0 0 0 1 1 *；
            `@monthly`：表示每月第一天的 0 点。等价于 0 0 0 1 * *；
            `@weekly`：表示每周第一天的 0 点，注意第一天为周日，即周六结束，周日开始的那个 0 点。等价于 0 0 0 * * 0；
            `@daily`：也可以写作@midnight，表示每天 0 点。等价于 0 0 0 * * *；
            @hourly：表示每小时的开始。等价于 0 0 * * * *。

        手动的固定时间间隔
            @every <duration>，手动设定时间间隔，只有 h m s （时分秒）可用，如：
                @every 1s，表示每秒，等价于 " * * * * * ? "
                @every 1h30m10s，表示每 1 小时 30 分 10 秒

*/
