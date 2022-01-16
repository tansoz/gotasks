package gotasks

type Worker interface {
	Quit() error
}

type worker struct {
	signal chan bool
	queue  <-chan func()
}

func NewWorker(queue <-chan func(), logfn func(...interface{})) Worker {
	w := new(worker)

	w.queue = queue
	w.signal = make(chan bool, 1)

	go w.run(logfn)

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
		case task := <-w.queue:
			task()
		}
	}
}
func (w *worker) close() {
	close(w.signal)
}
