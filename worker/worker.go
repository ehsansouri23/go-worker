package worker

import "sync"

type Worker[I any, O any] struct {
	wg         sync.WaitGroup
	jobChan    chan I
	resultChan chan O
}

func NewWorker[I any, O any]() *Worker[I, O] {
	return &Worker[I, O]{
		jobChan:    make(chan I),
		resultChan: make(chan O),
	}
}

func (w *Worker[I, O]) GetResults() chan O {
	return w.resultChan
}

func (w *Worker[I, O]) Wait() {
	w.wg.Wait()
	close(w.resultChan)
	close(w.jobChan)
}

func (w *Worker[I, O]) SubmitJob(job I) {
	w.wg.Add(1)
	w.jobChan <- job
}

func (w *Worker[I, O]) HandleJobs(function func(I) O) {
	for job := range w.jobChan {
		go func(job I) {
			w.resultChan <- function(job)
			w.wg.Done()
		}(job)
	}
}
