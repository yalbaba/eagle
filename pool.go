package eagle

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type sig struct{}

type Pool struct {
	// 协程池的最大容量
	capacity int32

	// 正在运行的协程数
	running int32

	// 空闲的worker超时时间，表示一直空闲很长时间的worker，超过时间被清理掉
	expiryDuration time.Duration

	// worker协程池处理核心对象,workers为池里的空闲worker集合
	workers []*Worker

	// 协程池关闭信号
	release chan sig

	// 保证数据处理安全
	lock sync.Mutex

	// once用在确保 Pool 关闭操作只会执行一次
	once sync.Once
}

func NewPool(opts ...option) *Pool {
	conf := &Options{}
	for _, o := range opts {
		o(conf)
	}
	if conf.capacity == 0 {
		conf.capacity = 10
	}
	if conf.expiryDuration == 0 {
		conf.expiryDuration = 300
	}
	p := &Pool{
		capacity:       conf.capacity,
		expiryDuration: conf.expiryDuration,
		workers:        make([]*Worker),
		release:        make(chan sig, 1),
	}
	// 开启线程清除空闲太久的worker
	p.monitorAndClear()
	return p, nil
}

// 提交任务到worker
func (p *Pool) Submit(task f) error {
	if len(p.release) > 0 {
		return fmt.Errorf("该协程池已经关闭")
	}
	w := p.getWorker()
	// 将任务放入woker里
	w.task <- task
	return nil
}

// 回收处理完任务的worker
func (p *Pool) putWorker(worker *Worker) {
	// 记录回收时间
	worker.recycle = time.Now()
	p.lock.Lock()
	// 放入协程池的worker队列中
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
}

// 获取可用的worker
func (p *Pool) getWorker() *Worker {
	var w *Worker
	var waiting bool
	p.lock.Lock()
	// 计算空闲的worker数量
	num := len(p.workers) - 1
	// 判断是否还有空闲的worker
	if num < 0 {
		// 没有空闲的worker,判断是否需要等待worker（当前正在运行的worke是否超过池最大容量）
		waiting = p.Running() >= p.Cap()
	} else {
		// 有空闲的worker，取出一个worker
		w = p.workers[num]
		p.workers = p.workers[:num]
	}
	p.lock.Unlock()

	if waiting {
		// 表示没有空闲的worke且正在运行的worker数量等于最大容量
		// 循环阻塞等待空闲中的worker
		for {
			p.lock.Lock()
			// 判断是否有空闲worker
			num = len(p.workers) - 1
			if num < 0 {
				p.lock.Unlock()
				continue
			}
			// 取出空闲的worker
			w = p.workers[num]
			p.workers = p.workers[:num]
			p.lock.Unlock()
			break
		}
	} else if w == nil {
		// 表示没有空闲worker但pool没有超过容量，则开启一个worker处理任务
		w = &Worker{
			pool: p,
			task: make(chan f, 1),
		}
		// 监听任务通道，处理任务
		w.run()
		// 池的运行worker数量加1
		p.incryRunning()
	}

	// 有空闲worker且容量也没满，则直接返回出去给client放入任务
	return w
}

// 并发安全获取正在运行的协程数
func (p *Pool) Running() int32 {
	// atomic.Load*系列函数只能保证读取的不是正在写入的值
	return atomic.LoadInt32(&p.running)
}

// 并发安全获取协程池的最大容量
func (p *Pool) Cap() int32 {
	// atomic.Load*系列函数只能保证读取的不是正在写入的值
	return atomic.LoadInt32(&p.capacity)
}

func (p *Pool) incryRunning() {
	atomic.AddInt32(&p.running, 1)
}
