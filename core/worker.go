package core

import (
	"context"
	"sync"
	"sync/atomic"
)

// WorkerPool Worker Pool
type WorkerPool struct {
	taskQueue  chan Task
	workers    []*Worker
	minWorkers int
	maxWorkers int
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	running    int32
	active     int32
}

// Task 任务
type Task struct {
	Handler func()
	Ctx     context.Context
}

// Worker Worker
type Worker struct {
	pool *WorkerPool
	id   int
}

// NewWorkerPool 创建Worker Pool
func NewWorkerPool(minWorkers, maxWorkers, queueSize int) *WorkerPool {
	if minWorkers <= 0 {
		minWorkers = 100
	}
	if maxWorkers <= 0 {
		maxWorkers = 10000
	}
	if queueSize <= 0 {
		queueSize = 100000
	}

	ctx, cancel := context.WithCancel(context.Background())
	pool := &WorkerPool{
		taskQueue:  make(chan Task, queueSize),
		workers:    make([]*Worker, 0, maxWorkers),
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动最小数量的Worker
	for i := 0; i < minWorkers; i++ {
		pool.addWorker(i)
	}

	return pool
}

// addWorker 添加Worker
func (p *WorkerPool) addWorker(id int) *Worker {
	worker := &Worker{
		pool: p,
		id:   id,
	}
	p.workers = append(p.workers, worker)

	p.wg.Add(1)
	go worker.start()

	return worker
}

// Submit 提交任务
func (p *WorkerPool) Submit(task Task) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case p.taskQueue <- task:
		// 动态调整Worker数量
		p.adjustWorkers()
		return nil
	}
}

// SubmitAsync 异步提交任务
func (p *WorkerPool) SubmitAsync(task Task) {
	select {
	case <-p.ctx.Done():
		return
	case p.taskQueue <- task:
		// 动态调整Worker数量
		p.adjustWorkers()
	}
}

// adjustWorkers 动态调整Worker数量
func (p *WorkerPool) adjustWorkers() {
	// 计算队列使用率
	queueLen := len(p.taskQueue)
	queueCap := cap(p.taskQueue)
	usage := float64(queueLen) / float64(queueCap)

	// 获取当前活跃Worker数
	activeWorkers := atomic.LoadInt32(&p.active)
	currentWorkers := len(p.workers)

	// 如果队列使用率超过80%且当前Worker数小于最大值，增加Worker
	if usage > 0.8 && currentWorkers < p.maxWorkers {
		newWorkers := min(currentWorkers*2, p.maxWorkers)
		for i := currentWorkers; i < newWorkers; i++ {
			p.addWorker(i)
		}
	}

	// 如果队列使用率低于20%且当前Worker数大于最小值，减少Worker
	if usage < 0.2 && currentWorkers > p.minWorkers && activeWorkers < int32(currentWorkers/2) {
		// 通过取消context来停止多余的Worker
		// 这里简化处理，实际实现可能需要更复杂的逻辑
	}
}

// Stop 停止Worker Pool
func (p *WorkerPool) Stop() {
	p.cancel()
	close(p.taskQueue)
	p.wg.Wait()
}

// Stats 获取统计信息
func (p *WorkerPool) Stats() WorkerPoolStats {
	return WorkerPoolStats{
		ActiveWorkers: atomic.LoadInt32(&p.active),
		TotalWorkers:  int32(len(p.workers)),
		QueueLength:   int32(len(p.taskQueue)),
		QueueCapacity: int32(cap(p.taskQueue)),
	}
}

// WorkerPoolStats Worker Pool统计信息
type WorkerPoolStats struct {
	ActiveWorkers int32
	TotalWorkers  int32
	QueueLength   int32
	QueueCapacity int32
}

// start 启动Worker
func (w *Worker) start() {
	defer w.pool.wg.Done()

	for {
		select {
		case <-w.pool.ctx.Done():
			return
		case task, ok := <-w.pool.taskQueue:
			if !ok {
				return
			}

			atomic.AddInt32(&w.pool.active, 1)
			task.Handler()
			atomic.AddInt32(&w.pool.active, -1)
		}
	}
}

// min 获取最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
