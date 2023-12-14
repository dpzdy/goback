package test

import (
	"fmt"
	"sync"
	"testing"
)

// Worker 表示工作池中的工作人员。
type Worker struct {
	ID int
}

type Pool struct {
	Workers   []Worker
	TaskQueue chan func()
	wg        sync.WaitGroup
}

func NewPool(numWorkers, taskSize int) *Pool {
	pool := &Pool{TaskQueue: make(chan func(), taskSize)}
	pool.Workers = make([]Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		pool.Workers[i] = Worker{
			ID: i,
		}
	}
	return pool
}
func (p *Pool) Submit(task func()) {
	p.wg.Add(1)
	p.TaskQueue <- task
}
func (p *Pool) ShutDown() {
	close(p.TaskQueue)
	p.wg.Wait()
}

func (w Worker) Start(p *Pool) {
	go func() {
		for task := range p.TaskQueue {
			task() // 执行任务函数
			p.wg.Done()
		}
	}()
}
func TestPool(t *testing.T) {
	pool := NewPool(3, 10)

	// 启动工作人员
	for _, worker := range pool.Workers {
		worker.Start(pool)
	}

	// 提交一些任务到线程池
	for i := 0; i < 100; i++ {
		id := i
		pool.Submit(func() {
			fmt.Printf("Task %d executed by Worker %d\n", id, id%3)
		})
	}

	// 关闭线程池并等待所有任务完成
	pool.ShutDown()

}
