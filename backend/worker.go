package main

import astiworker "github.com/asticode/go-astitools/worker"

type worker struct {
	s *session
	w *astiworker.Worker
}

func newWorker() *worker {
	return &worker{
		w: astiworker.NewWorker(),
	}
}

func (w *worker) wait() {
	// Handle signals
	w.w.HandleSignals()

	// Wait
	w.w.Wait()
}
