package main

import (
	"github.com/asticode/go-asting"
	astiworker "github.com/asticode/go-astitools/worker"
	"go.etcd.io/bbolt"
)

type worker struct {
	db *bbolt.DB
	i  *asting.Client
	w  *astiworker.Worker
}

func newWorker(i *asting.Client) *worker {
	return &worker{
		i: i,
		w: astiworker.NewWorker(),
	}
}

func (w *worker) wait() {
	// Handle signals
	w.w.HandleSignals()

	// Wait
	w.w.Wait()
}
