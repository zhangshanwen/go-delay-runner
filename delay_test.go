package go_delay_runner

import (
	"fmt"
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
				ExecuteAt: time.Now().Add(time.Duration(i+2) * time.Second),
				Handlers: []Handler{func(...interface{}) {
					fmt.Println("开始执行任务", i, time.Now())
				}},
			}
			w.Push(t)
		}(i)
	}
	wg.Wait()
	fmt.Println("-----")
	for {
		time.Sleep(3 * time.Second)
	}
}
