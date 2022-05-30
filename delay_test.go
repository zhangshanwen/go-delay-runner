package go_delay_runner

import (
	"sync"
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	w := NewWorker()
	nums := 10
	wg := sync.WaitGroup{}
	wg.Add(nums)
	for i := 0; i < nums; i++ {
		go func(i int) {
			defer wg.Done()
			task := Task{
				ExecuteAt: time.Now().Add(time.Duration(i) * time.Second),
				Handlers: []Handler{func(interface{}) {
					t.Log("开始执行任务", i, time.Now())
				}},
			}
			t.Log(w.Push(&task))
			t.Log(w.Push(&task))
			t.Log(w.Push(&task))
			t.Log(w.Push(&task))
		}(i)
	}
	wg.Wait()
	t.Log("-----")
	for {
		time.Sleep(3 * time.Second)
	}
}

func TestDelayWork(t *testing.T) {
	w := NewWorker()
	t.Log("创建时间", time.Now())
	type A struct {
		Name string
	}
	a := A{Name: "dfss"}
	task := Task{
		ExecuteAt: time.Now().Add(time.Duration(2) * time.Second),
		Args:      a,
	}
	task.Handlers = []Handler{func(a interface{}) {
		t.Log("开始执行任务", 1, task.ExecuteAt, time.Now(), a)
	}}
	w.Push(&task)
	time.Sleep(3 * time.Second)
	b := A{Name: "lisi"}
	task = Task{
		ExecuteAt: time.Now().Add(time.Duration(2) * time.Second),
		Args:      b,
	}
	task.Handlers = []Handler{func(a interface{}) {
		if p, ok := a.(A); ok {
			t.Log("开始执行任务", 2, time.Now(), p.Name)
		} else {
			t.Log("开始执行任务", 2, time.Now(), "执行失败")
		}

	}}
	w.Push(&task)
	t.Log("-----")
	for {
		time.Sleep(3 * time.Second)
	}
}
