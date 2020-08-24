# Eagle 
示例：
1、构建协程池
p := eagle.NewPool(参数选填，参考example)

2、执行的job需要实现Handler接口的process方法，参数由实现接口的结构体传递
type TestHandler struct {
	id int
}

func (t *TestHandler) Process() {
	fmt.Println("process job::",t.id)
	time.Sleep(1 * time.Second)
}

3、使用submit提交任务
p.Submit(t)

