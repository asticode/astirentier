package main

import (
	astiworker "github.com/asticode/go-astitools/worker"
	"go.etcd.io/bbolt"
)

type worker struct {
	db *bbolt.DB
	w  *astiworker.Worker
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
