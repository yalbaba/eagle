package eagle

import "time"

// task方法需自定义
type f func() error

type Worker struct {
	// 属于哪个池
	pool *Pool

	// 存任务方法的管道
	task chan f

	// 该worker重新放入队列的时间
	recycle time.Time
}

// worker 处理每一个任务
func (n *Worker) run() {
	// 开启一个协程处理放入管道的任务
	go func() {
		for v := range n.task {
			if v == nil {
				// 协程池的正在运行worker数减一
			}
			// 执行传入的任务
			v()
			// 执行完后将worker放入协程池的worker集合
			n.pool.putWorker(n)
		}
	}()
}
