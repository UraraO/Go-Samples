package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
)

// 使用cron，参考README.go

// ////////////////////////////////

// 暂定类型：
/*
	本地：（不考虑实例问题，本地互斥由用户在Run中自主加锁实现
		BasicCron，CronWithLog
	分布式单实例：
		数据不同步，分布式锁应当与CronManager无关，而是某个分布式系统的全局变量，在创建CronManager时，将分布式锁作为参数传给Cron
*/

// Log,暂定使用xfx库中的el日志，具体实现方案：
/*
	1. 定时任务日志使用el输出到stdout
	2. 业务日志交由用户实现，用户将日志写在Run中
*/

// 记录任务的调用，方案
/*
	持久化：xfx的日志el
*/

// 任务完成时通知，方案
/*
	用户设定一个回调方法，Run结束时调用
*/

type TaskDescription struct {
	ID          cron.EntryID
	Description string
	Skip        bool
	CoverPanic  bool
}

type CallbackFunc interface {
	Callback()
}

type TaskWithCallback struct {
	Task     cron.Job
	Callback CallbackFunc
}

func (thisTask *TaskWithCallback) Run() {
	thisTask.Task.Run()
	thisTask.Callback.Callback()
}

type BasicCron struct {
	Cron  *cron.Cron                       // 定时任务管理器
	Tasks map[cron.EntryID]TaskDescription // 正在运行的任务
	// 					ID - 任务描述
	Running bool

	// 成员方法
	//
	// 管理任务
	// AddFunc
	// AddTask	TODO task的Run方法为指针接收器时，传入&task
	// Remove	根据ListExistTask的结果，或者如果当前Cron仅有一个任务，直接Stop
	//
	// 启动Cron
	// Start
	// Stop
	//
	// 特殊操作
	// ListExistTask	列出当前Cron管理器管理的所有任务ID及其描述，方便根据ID进行任务删除
	//
}

func NewBasicCron() *BasicCron {
	bcron := &BasicCron{
		Cron:    cron.New(cron.WithSeconds()),
		Tasks:   make(map[cron.EntryID]TaskDescription, 64),
		Running: false,
	}
	return bcron
}

func (thisCron *BasicCron) AddFunc(cronSpec string, fun func(), funcDescription string) (cron.EntryID, error) {
	// TODO 添加日志el
	// task := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger))
	// taskID, err := thisCron.Cron.AddFunc(cronSpec, task.Then(fun))
	taskID, err := thisCron.Cron.AddFunc(cronSpec, fun)
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        true,
		CoverPanic:  false,
	}
	return taskID, err
}

func (thisCron *BasicCron) AddTask(cronSpec string, task cron.Job, funcDescription string) (cron.EntryID, error) {
	// TODO 添加日志el
	chain := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	taskID, err := thisCron.Cron.AddJob(cronSpec, chain.Then(task))
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        true,
		CoverPanic:  true,
	}
	return taskID, err
}

func (thisCron *BasicCron) AddTaskWithOption(cronSpec string, task cron.Job, funcDescription string, delay bool, coverPanic bool) (cron.EntryID, error) {
	// TODO 添加日志el
	var chain cron.Chain
	if delay && coverPanic {
		chain = cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	} else if !delay && coverPanic {
		chain = cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	} else if delay && !coverPanic {
		chain = cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger))
	} else if !delay && !coverPanic {
		chain = cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger))
	}

	taskID, err := thisCron.Cron.AddJob(cronSpec, chain.Then(task))
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        !delay,
		CoverPanic:  coverPanic,
	}
	return taskID, err
}

func (thisCron *BasicCron) AddTaskWithCallback(cronSpec string, task cron.Job, funcDescription string, callback CallbackFunc) (cron.EntryID, error) {
	// TODO 添加日志el
	ctask := TaskWithCallback{
		Task:     task,
		Callback: callback,
	}
	chain := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	taskID, err := thisCron.Cron.AddJob(cronSpec, chain.Then(&ctask))
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        true,
		CoverPanic:  true,
	}
	return taskID, err
}

func (thisCron *BasicCron) Remove(taskID cron.EntryID) {
	desc := thisCron.Tasks[taskID]
	thisCron.Cron.Remove(taskID)
	delete(thisCron.Tasks, taskID)
	fmt.Println("Remove task:", desc)
	// el.LOG(remove task: desc)
}

func (thisCron *BasicCron) ListExistTask() {
	for _, task := range thisCron.Tasks {
		fmt.Println(task)
	}
}

func (thisCron *BasicCron) Start() {
	thisCron.Cron.Start()
	thisCron.Running = true
	// el.LOG(cron start)
}

func (thisCron *BasicCron) Stop() context.Context {
	ctx := thisCron.Cron.Stop()
	thisCron.Running = false
	// el.LOG(cron stop)
	return ctx
}

// ///////////////////////////
//
// TODO 所有分布式锁操作

type DistMutex interface {
	Lock()
	Unlock()
}

type DistributeCron struct {
	Cron  *cron.Cron                       // 定时任务管理器
	Tasks map[cron.EntryID]TaskDescription // 正在运行的任务
	// 					ID - 任务描述
	DistMut DistMutex
	Running bool

	// 成员方法
	//
	// 管理任务
	// AddFunc
	// AddTask	TODO task的Run方法为指针接收器时，传入&task
	// Remove	根据ListExistTask的结果，或者如果当前Cron仅有一个任务，直接Stop
	//
	// 启动Cron
	// Start
	// Stop
	//
	// 特殊操作
	// ListExistTask	列出当前Cron管理器管理的所有任务ID及其描述，方便根据ID进行任务删除
	//
}

func NewDistributeCron(mutex DistMutex) *DistributeCron {
	dcron := &DistributeCron{
		Cron:    cron.New(cron.WithSeconds()),
		Tasks:   make(map[cron.EntryID]TaskDescription, 64),
		DistMut: mutex,
		Running: false,
	}
	return dcron
}

func (thisCron *DistributeCron) AddFunc(cronSpec string, fun func(), funcDescription string) (cron.EntryID, error) {
	thisCron.DistMut.Lock()
	defer thisCron.DistMut.Unlock()
	// TODO 添加日志el
	taskID, err := thisCron.Cron.AddFunc(cronSpec, fun)
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        true,
		CoverPanic:  false,
	}
	return taskID, err
}

func (thisCron *DistributeCron) AddTask(cronSpec string, task cron.Job, funcDescription string) (cron.EntryID, error) {
	// TODO 添加日志el
	thisCron.DistMut.Lock()
	defer thisCron.DistMut.Unlock()
	chain := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	taskID, err := thisCron.Cron.AddJob(cronSpec, chain.Then(task))
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        true,
		CoverPanic:  true,
	}
	return taskID, err
}

func (thisCron *DistributeCron) AddTaskWithOption(cronSpec string, task cron.Job, funcDescription string, delay bool, coverPanic bool) (cron.EntryID, error) {
	// TODO 添加日志el
	thisCron.DistMut.Lock()
	defer thisCron.DistMut.Unlock()
	var chain cron.Chain
	if delay && coverPanic {
		chain = cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	} else if !delay && coverPanic {
		chain = cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	} else if delay && !coverPanic {
		chain = cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger))
	} else if !delay && !coverPanic {
		chain = cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger))
	}

	taskID, err := thisCron.Cron.AddJob(cronSpec, chain.Then(task))
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        !delay,
		CoverPanic:  coverPanic,
	}
	return taskID, err
}

func (thisCron *DistributeCron) AddTaskWithCallback(cronSpec string, task cron.Job, funcDescription string, callback CallbackFunc) (cron.EntryID, error) {
	// TODO 添加日志el
	thisCron.DistMut.Lock()
	defer thisCron.DistMut.Unlock()
	ctask := TaskWithCallback{
		Task:     task,
		Callback: callback,
	}
	chain := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger))
	taskID, err := thisCron.Cron.AddJob(cronSpec, chain.Then(&ctask))
	if err != nil {
		// el.ERROR(add job fail)
		return -1, err
	}
	// el.LOG(add job success)
	thisCron.Tasks[taskID] = TaskDescription{
		ID:          taskID,
		Description: funcDescription,
		Skip:        true,
		CoverPanic:  true,
	}
	return taskID, err
}

func (thisCron *DistributeCron) Remove(taskID cron.EntryID) {
	thisCron.DistMut.Lock()
	defer thisCron.DistMut.Unlock()
	desc := thisCron.Tasks[taskID]
	thisCron.Cron.Remove(taskID)
	delete(thisCron.Tasks, taskID)
	fmt.Println("Remove task:", desc)
	// el.LOG(remove task: desc)
}

func (thisCron *DistributeCron) ListExistTask() {
	thisCron.DistMut.Lock()
	defer thisCron.DistMut.Unlock()
	for _, task := range thisCron.Tasks {
		fmt.Println(task)
	}
}

// NOTE 分布式锁的使用
// NOTE 如果仅允许一个实例运行，则按下方实现，Start加锁，Stop解锁
// NOTE 如果允许多个实例交替执行，仅Run函数互斥，则应在Run前加锁，Run后解锁

func (thisCron *DistributeCron) Start() {
	thisCron.DistMut.Lock()
	thisCron.Cron.Start()
	thisCron.Running = true
	// el.LOG(cron start)
}

func (thisCron *DistributeCron) Stop() context.Context {
	defer thisCron.DistMut.Unlock()
	ctx := thisCron.Cron.Stop()
	thisCron.Running = false
	// el.LOG(cron stop)
	return ctx
}
