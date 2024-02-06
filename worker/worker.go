package worker

import "sync"

const (
	WithoutGoroutinePool = -1
)

type Handler[I any, O any] func(I) O

type Worker[I any, O any] struct {
	wg            sync.WaitGroup
	jobChan       chan I
	resultChan    chan O
	handler       Handler[I, O]
	maxGoroutines int
}

func NewWorker[I any, O any](maxGoroutines int) *Worker[I, O] {
	return &Worker[I, O]{
		jobChan:       make(chan I),
		resultChan:    make(chan O),
		maxGoroutines: maxGoroutines,
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
	w.jobChan <- job
}

func (w *Worker[I, O]) SetHandler(handler Handler[I, O]) {
	w.handler = handler
}

func (w *Worker[I, O]) workerGoroutine(id int) {
	defer w.wg.Done()
	for job := range w.jobChan {
		w.resultChan <- w.handler(job)
		w.wg.Done()
	}
}

func (w *Worker[I, O]) HandleJobs() {
	if w.maxGoroutines == WithoutGoroutinePool {
		for job := range w.jobChan {
			go func(job I) {
				w.resultChan <- w.handler(job)
				w.wg.Done()
			}(job)
		}
	} else {
		for i := 1; i <= w.maxGoroutines; i++ {
			w.wg.Add(1)
			go w.workerGoroutine(i)
		}
	}
}
