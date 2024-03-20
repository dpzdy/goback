package test

import (
	"fmt"
	"testing"
	"time"
)

type Task struct {
	f func() func() string
}

func NewTask(f func() func() string) *Task {
	task := &Task{
		f: f,
	}
	return task
}

func (t *Task) Execute() {
	cur := t.f()
	fmt.Println(cur())
}

type MyPool struct {
	EntryChannel chan *Task
	WorkerNum    int
	JobsChannel  chan *Task
}

func NewMyPool(cap int) *MyPool {
	pool := &MyPool{
		EntryChannel: make(chan *Task),
		WorkerNum:    cap,
		JobsChannel:  make(chan *Task),
	}
	return pool
}

func (p *MyPool) Worker(workerID int) {
	for task := range p.JobsChannel {
		task.Execute()
		fmt.Println("worker ID is ", workerID, " Done")
	}
}

func (p *MyPool) Run() {
	for i := 0; i < p.WorkerNum; i++ {
		fmt.Println("worker ID is ", i, " Begin")
		go p.Worker(i)
	}
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}
	close(p.JobsChannel)
	fmt.Println("Close JobsChannel")
	close(p.EntryChannel)
	fmt.Println("Close EntryChannel")
	//println("Close EntryChannel")
}

func TestMyPool(t *testing.T) {

	task := NewTask(func() func() string {
		//time.Sleep(time.Second)

		return func() string {
			return fmt.Sprintf("创建一个Task: %s", time.Now().Format("2006-01-02 15:04:05"))
		}

	})
	p := NewMyPool(3)
	go func() {
		triker := time.NewTicker(5 * time.Second)
		for {
			p.EntryChannel <- task
		}
		go func() {
			select {
			//case p.EntryChannel <- task:

			case <-triker.C:
				break
				//default:
				//	p.EntryChannel <- task
			}
		}()
	}()

	p.Run()
}
