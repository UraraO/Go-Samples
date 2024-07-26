package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

type somejob struct {
	count int
}

func (job *somejob) Run() {
	job.count++
	fmt.Println(job.count)
}

func TestBasicCron(t *testing.T) {
	// 创建Task
	job1 := somejob{
		count: 200,
	}
	job2 := somejob{
		count: 30000,
	}
	// 创建Cron
	bc := NewBasicCron()
	// 添加三个Cron
	ID, err := bc.AddFunc("@every 1s", func() {
		fmt.Println("func 1s do")
	}, "func 1s")
	if err != nil {
		fmt.Println(err)
		t.Fatal("AddFunc err")
		return
	}
	fmt.Println("func ID:", ID)

	taskID1, err := bc.AddTask("@every 2s", &job1, "Task 2s")
	if err != nil {
		t.Fatal("AddTask err")
		return
	}
	fmt.Println("Task ID:", taskID1)

	taskID2, err := bc.AddTaskWithOption("@every 3s", &job2, "Task 3s", true, true)
	if err != nil {
		t.Fatal("AddTask err")
		return
	}
	fmt.Println("Task ID:", taskID2)
	fmt.Println()
	// 运行，测试列表和删除任务功能
	bc.ListExistTask()
	bc.Start()
	time.Sleep(10 * time.Second)
	bc.Stop()
	bc.Remove(ID)
	bc.Remove(taskID1)
	fmt.Println()
	bc.ListExistTask()
	bc.Start()
	time.Sleep(5 * time.Second)
	bc.Stop()
	fmt.Println()
	bc.ListExistTask()
	bc.Remove(taskID2)
	fmt.Println()
	bc.ListExistTask()
}

type callbackFunc struct {
}

func (callback *callbackFunc) Callback() {
	fmt.Println("Callbackkkkkkkkkkkkkkkkkkk")
}

func TestBasicCronWithCallback(t *testing.T) {
	// 创建Task
	job1 := somejob{
		count: 0,
	}
	cb := callbackFunc{}
	// 创建Cron
	bc := NewBasicCron()

	taskID1, err := bc.AddTaskWithCallback("0/3 * * * * ?", &job1, "Task 3s_1", &cb)
	if err != nil {
		t.Fatal("AddTask err")
		return
	}
	fmt.Println("Task ID:", taskID1)

	// 运行，测试列表和删除任务功能
	bc.ListExistTask()
	bc.Start()
	time.Sleep(10 * time.Second)
	bc.Stop()
}

// discard all content below
func SampleBasicCronWithSeconds() {
	cs := cron.New() // 默认不带秒级控制
	// cs := cron.New(cron.WithSeconds()) // 使用WithSeconds赋予秒级控制

	// cs.AddFunc("* * * * * ?", func() { fmt.Println("cs,\"*****?\" :", time.Now()) }) // cron表达式
	cs.AddFunc("@every 1s", func() { fmt.Println("cs,\"every 1s\" :", time.Now()) }) // 预定义的方便用法

	// cs.Start() 	// 异步运行，新开goroutine执行定时任务
	cs.Run() // 同步运行
}

type SomeTask struct {
	tm time.Time
}

func (tsk *SomeTask) Run() {
	tsk.tm = time.Now()
	fmt.Println(tsk.tm)
}

type panicJob struct {
	count int
}

func (p *panicJob) Run() {
	p.count++
	if p.count == 1 {
		panic("oooooooooooooops!!!")
	}

	fmt.Println("hello world")
}

func TestBasicUsage(t *testing.T) {
	c := cron.New()
	_, err := c.AddFunc("@every 1s", func() { fmt.Println("hello") })
	if err != nil {
		return
	}

	// cs.Start() 	// 异步运行，新开goroutine执行定时任务
	c.Run() // 同步运行
	_ = c.Stop()

}
