package cron

import (
	"sync"
	"sync/atomic"
	"time"
)

type ITask interface {
	// 返回任务名
	Name() string
	// 获取handler
	Handler() Handler
	// 返回启用状态
	IsEnable() bool

	// 获取触发时间
	TriggerTime() time.Time
	// 生成下次触发时间, 如果返回了 false 表示没有下一次了, 返回的时间一定>t
	MakeNextTriggerTime(t time.Time) (time.Time, bool)
	// 立即触发执行, 阻塞等待执行结束
	Trigger(ctx IContext, callback ErrCallback) error

	// 重置定时, 发生在被定时器添加任务时和重新设为启用时
	resetClock()
	// 设置启用
	setEnable(enable bool)
	// 设置堆索引
	setHeapIndex(index int)
	// 获取堆索引
	getHeapIndex() int
}

type Task struct {
	name string

	handler Handler

	triggerTime time.Time
	trigger     ITrigger
	executor    IExecutor

	enable int32
	mx     sync.Mutex // 用于锁 triggerTime, trigger, executor

	heapIndex int // 堆索引
}

type TaskConfig struct {
	Trigger  ITrigger
	Executor IExecutor
	Handler  Handler
	Enable   bool
}

// 创建一个任务
func NewTask(name string, expression string, enable bool, handler Handler) ITask {
	trigger := NewCronTrigger(expression)
	executor := NewExecutor(0, 0, 1)
	return NewTaskOfConfig(name, TaskConfig{
		Trigger:  trigger,
		Executor: executor,
		Handler:  handler,
		Enable:   enable,
	})
}

// 根据任务配置创建一个任务
func NewTaskOfConfig(name string, config TaskConfig) ITask {
	t := &Task{
		name:     name,
		trigger:  config.Trigger,
		executor: config.Executor,
		handler:  config.Handler,
	}
	t.setEnable(config.Enable)
	return t
}

func (t *Task) Name() string {
	return t.name
}
func (t *Task) Handler() Handler {
	return t.handler
}
func (t *Task) IsEnable() bool {
	return atomic.LoadInt32(&t.enable) == 1
}
func (t *Task) TriggerTime() time.Time {
	t.mx.Lock()
	tt := t.triggerTime
	t.mx.Unlock()
	return tt
}
func (t *Task) MakeNextTriggerTime(tt time.Time) (time.Time, bool) {
	if !t.IsEnable() {
		return tt, false
	}

	tt, ok := t.trigger.MakeNextTriggerTime(tt)
	t.mx.Lock()
	t.triggerTime = tt
	t.mx.Unlock()
	return tt, ok
}
func (t *Task) Trigger(ctx IContext, callback ErrCallback) error {
	return t.execute(ctx, callback)
}

// 执行
func (t *Task) execute(ctx IContext, errCallback ErrCallback) error {
	t.mx.Lock()
	executor := t.executor
	t.mx.Unlock()
	return executor.Do(ctx, errCallback)
}

func (t *Task) resetClock() {
	t.mx.Lock()
	trigger := t.trigger
	t.mx.Unlock()
	trigger.ResetClock()
}
func (t *Task) setEnable(enable bool) {
	if enable {
		atomic.StoreInt32(&t.enable, 1)
	} else {
		atomic.StoreInt32(&t.enable, 0)
	}
}
func (t *Task) setHeapIndex(index int) {
	t.heapIndex = index
}
func (t *Task) getHeapIndex() int {
	return t.heapIndex
}
