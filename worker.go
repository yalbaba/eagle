package eagle

import "time"

// task方法需自定义（没想好如何实现更通用）
type handle func(...interface{}) error

type Worker struct {
	// 属于哪个池
	pool *Pool

	// 存任务方法的管道
	task chan handle

	// 该worker重新放入队列的时间
	recycleTime time.Time
}

// worker 处理每一个任务
func (n *Worker) run() {
	// 开启一个协程处理放入管道的任务
	go func() {
		for f := range n.task {
			if f == nil {
				// 协程池的正在运行worker数减一
				n.pool.dcrRunning()
				return
			}
			// 执行传入的任务
			f()
			// 执行完后将worker放入协程池的worker集合
			n.pool.putWorker(n)
		}
	}()
}
