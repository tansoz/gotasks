package gotasks

import (
	"errors"
	"log"
)

type Tasks interface {
	Async(func()) error
	Sync(func()) error
	Quit() error
	Active() int
}

type tasks struct {
	queue   chan func()
	signal  chan bool
	logfn   func(...interface{})
	workers []Worker
	active  int
}

func NewTasks(workers int, fn func(Tasks), logfn func(...interface{})) error {
	tks := new(tasks)

	tks.queue = make(chan func(), 20)
	tks.signal = make(chan bool, 1)
	tks.workers = make([]Worker, workers)
	tks.logfn = logfn
	tks.active = 0

	if tks.logfn == nil {
		tks.logfn = log.Println
	}

	tks.newWorkers(workers)

	go fn(tks)

	return nil
}
func (tks *tasks) close() {
	close(tks.queue)
	close(tks.signal)
}
func (tks *tasks) newWorkers(workers int) error {
	for i := 0; i < workers; i++ {
		tks.workers[i] = tks.NewWorker()
	}
	return nil
}
func (tks *tasks) Async(task func()) error {
	if task == nil {
		return errors.New("the task can not be a null")
	}
	tks.queue <- task
	return nil
}
func (tks *tasks) Sync(task func()) error {
	if task == nil {
		return errors.New("the task can not be a null")
	}
	signal := make(chan bool, 1) // 任务完成信号
	tks.queue <- func() {
		task()
		signal <- true
	}
	<-signal
	return nil
}
func (tks *tasks) Active() int {
	return tks.active
}
func (tks *tasks) Quit() error {
	defer tks.quitWorkers()
	defer tks.close()
	tks.signal <- true
	return nil
}
func (tks *tasks) quitWorkers() {
	for _, wk := range tks.workers {
		wk.Quit()
	}
}
