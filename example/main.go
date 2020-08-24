package main

import (
	"fmt"
	"time"

	eagle "yalbaba/eagle"
)

type TestHandler struct {
	id int
}

func (t *TestHandler) Process() {
	fmt.Println("process job::", t.id)
	time.Sleep(1 * time.Second)
}

func main() {
	p := eagle.NewPool(eagle.WithCapacity(50))

	now := time.Now()
	for i := 0; i < 1000; i++ {
		t := &TestHandler{
			id: i,
		}
		p.Submit(t)
	}
	cur := time.Now()
	fmt.Println("sub::::", cur.Sub(now))
}
