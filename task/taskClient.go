package task

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runer struct {
	Interrupt chan os.Signal   //发送的信号。从系统
	Complte   chan error       //通道报告处理任务已经完成
	Timeout   <-chan time.Time //任务超时
	Tasks     []func(interface{})
}

//定义error类型

var ErrTimeOut = errors.New("received timeout")
var ErrInterrupt = errors.New("received intterrupt")

//返回一个准备使用的Runer
func New(d time.Duration) *Runer {
	return &Runer{
		Interrupt: make(chan os.Signal, 1),
		Complte:   make(chan error),
		Timeout:   time.After(d),
	}
}

//添加任务到runer

func (r *Runer) Add(tasks ...func(interface{})) {
	r.Tasks = append(r.Tasks, tasks...)
}

//执行任务并监听通道

func (r *Runer) Start() error {
	signal.Notify(r.Interrupt, os.Interrupt)

	go func() {
		r.Complte <- r.run()
	}()

	select {
	case err := <-r.Complte:
		return err
	case <-r.Timeout:
		return ErrTimeOut
	}
}

//执行
func (r *Runer) run() error {
	for id, task := range r.Tasks {
		if r.gotInterrupt() {
			return ErrInterrupt
		}
		task(id) //func 参数
	}
	return nil
}

//判断中断信号
func (r *Runer) gotInterrupt() bool {
	select {
	case <-r.Interrupt:
		signal.Stop(r.Interrupt)
		return true
	default:
		return false

	}
}
