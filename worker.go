package gotasks

type Worker interface {
	Quit() error
}

type worker struct {
	signal chan bool
	tasks  *tasks
}

func (tks *tasks) NewWorker() Worker {
	w := new(worker)

	w.tasks = tks
	w.signal = make(chan bool, 1)

	go w.run(w.tasks.logfn)

	return w
}
func (w *worker) Quit() error {
	defer w.close()
	w.signal <- true
	return nil
}
func (w *worker) run(logfn func(...interface{})) {
	defer func() {
		if err := recover(); err != nil {
			w.close()
			logfn(err)
		}
	}()
	for {
		select {
		case <-w.signal:
			return
		case task := <-w.tasks.queue:
			w.tasks.active += 1
			task()
			w.tasks.active -= 1
		}
	}
}
func (w *worker) close() {
	close(w.signal)
}
