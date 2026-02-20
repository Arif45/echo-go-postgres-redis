package worker

import (
	"log"
	"time"
)

type Worker struct {
	interval time.Duration
	job      func() error
	name     string
}

func NewWorker(interval time.Duration, job func() error, name string) *Worker {
	return &Worker{
		interval: interval,
		job:      job,
		name:     name,
	}
}

func (w *Worker) Start() {
	log.Printf("[%s] worker started (interval: %v)", w.name, w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.run()
	for range ticker.C {
		w.run()
	}

}

func (w *Worker) run() {
	log.Printf("[worker] %s running...", w.name)
	if err := w.job(); err != nil {
		log.Printf("[worker] %s error: %v", w.name, err)
	} else {
		log.Printf("[worker] %s completed successfully", w.name)
	}
}

func (w *Worker) Stop() {
	// Implement the logic to gracefully stop the worker
}
