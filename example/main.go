package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	eagle "yalbaba/Eagle"
)

func main() {
	p := eagle.NewPool()
	t := func(ps ...int) error {
		fmt.Println("11")
		return nil
	}
	now := time.Now()
	for i := 0; i < 1000; i++ {
		p.Submit(t)
	}
	cur := time.Now()

	fmt.Println("sub::::", now.Sub(cur))

	time.Sleep(time.Second)
	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println("线程池启动了")
	<-c
	fmt.Println("结束")
}
