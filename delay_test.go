package go_delay_runner

import (
	"sync"
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	w := NewWorker()
	w.Start()
	nums := 100
	wg := sync.WaitGroup{}
	wg.Add(nums)
	for i := 0; i < nums; i++ {
		go func(i int) {
			defer wg.Done()
			t := Task{
				ExecuteAt: time.Now().Add(time.Duration(i) * time.Second),
				Handlers: []Handler{func(interface{}) {
					t.Log("开始执行任务", i, time.Now())
				}},
			}
			w.Push(&t)
			w.Push(&t)
			w.Push(&t)
		}(i + 5)
	}
	wg.Wait()
	t.Log("-----")
	for {
		time.Sleep(3 * time.Second)
	}
}

func TestDelayWork(t *testing.T) {
	w := NewWorker()
	w.Start()
	t.Log("创建时间", time.Now())
	a := struct {
		Name string
	}{Name: "zhangsan"}
	task := Task{
		ExecuteAt: time.Now().Add(time.Duration(2) * time.Second),
		Args:      a,
	}
	task.Handlers = []Handler{func(interface{}) {
		t.Log("开始执行任务", 1, time.Now(), task.Args)
	}}
	w.Push(&task)
	time.Sleep(3 * time.Second)
	b := struct {
		Name string
	}{Name: "lisi"}
	task = Task{
		ExecuteAt: time.Now().Add(time.Duration(2) * time.Second),
		Args:      b,
	}
	task.Handlers = []Handler{func(interface{}) {
		t.Log("开始执行任务", 2, time.Now(), task.Args)
	}}
	w.Push(&task)
	t.Log("-----")
	for {
		time.Sleep(3 * time.Second)
	}
}
