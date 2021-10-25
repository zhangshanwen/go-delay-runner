package go_delay_runner

import (
	"log"
	"sync"
	"time"
)

type (
	workSignalType int
	Handler        func(...interface{})
	Task           struct {
		ExecuteAt time.Time
		Handlers  []Handler
	}
	Logger interface {
		Println(v ...interface{})
	}
	Node struct {
		ExecuteAt time.Time
		NextNode  *Node
		Tasks     []Task
	}
	Worker struct {
		Mx          sync.Mutex
		NextNode    *Node
		Signal      chan workSignalType
		SignalStart chan workSignalType
		Logger      Logger
	}
)

const (
	workSignalStop workSignalType = iota + 1
	workSignalStart
)

var (
	DefaultCacheLen = 100
)

// Push 推送新任务
func (n *Node) Push(task Task, w *Worker) {
	for {
		if n.ExecuteAt.Equal(task.ExecuteAt) {
			// 如果当前任务跟即将执行任务时间相同
			n.Tasks = append(n.Tasks, task)
			return
		} else if n.ExecuteAt.After(task.ExecuteAt) {
			// 如果当前任务比即将执行任务时间更小
			oldN := *n
			n.NextNode = &oldN
			n.Tasks = []Task{task}
			n.ExecuteAt = task.ExecuteAt
			w.Signal <- workSignalStop
			go w.Run()

			return
		} else {
			// 如果当前任务比即将执行任务时间更大
			if n.NextNode == nil {
				n.NextNode = &Node{
					NextNode:  nil,
					Tasks:     []Task{task},
					ExecuteAt: task.ExecuteAt,
				}
				return
			}

			n.NextNode.Push(task, w)
			return
		}
	}

}

func (w *Worker) Push(task Task) {
	w.Mx.Lock()
	defer w.Mx.Unlock()
	if w.NextNode == nil {
		w.NextNode = &Node{
			NextNode:  nil,
			Tasks:     []Task{task},
			ExecuteAt: task.ExecuteAt,
		}
		w.SignalStart <- workSignalStart
		return
	}
	w.NextNode.Push(task, w)
}

func NewWorker() *Worker {
	return &Worker{
		Signal:      make(chan workSignalType, DefaultCacheLen),
		SignalStart: make(chan workSignalType, DefaultCacheLen),
		Logger:      log.Default(),
	}
}

func (w *Worker) Start() {
	go func() {
		for _ = range w.SignalStart {
			go w.Run()
		}
	}()
}

func (w *Worker) Run() {
	w.Mx.Lock()
	defer w.Mx.Unlock()
	if w.NextNode == nil {
		return
	}
	select {
	case s := <-w.Signal:
		switch s {
		case workSignalStop:
			w.Logger.Println("程序终止中.......")
			break
		}
	case <-time.After(w.NextNode.ExecuteAt.Sub(time.Now())):
		for _, task := range w.NextNode.Tasks {
			for _, handler := range task.Handlers {
				go handler()
			}
		}
		w.NextNode = w.NextNode.NextNode
		go w.Run()
	}
}
