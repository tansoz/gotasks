package gotasks

import "log"

type Tasks interface {
	Add(func()) error
	Quit() error
}

type tasks struct {
	queue   chan func()
	signal  chan bool
	logfn   func(...interface{})
	workers []Worker
}

func NewTasks(workers int, fn func(Tasks), logfn func(...interface{})) error {
	tks := new(tasks)

	tks.queue = make(chan func(), 20)
	tks.signal = make(chan bool, 1)
	tks.workers = make([]Worker, workers)
	tks.logfn = logfn

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
		tks.workers[i] = NewWorker(tks.queue, tks.logfn)
	}
	return nil
}
func (tks *tasks) Add(task func()) error {
	tks.queue <- task
	return nil
}
func (tks *tasks) Quit() error {
	defer tks.quitWorkers()
	defer tks.close()
	tks.signal <- true
	return nil
}
func (tks *tasks) quitWorkers() {
	for _, worker := range tks.workers {
		worker.Quit()
	}
}
