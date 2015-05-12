package replay

import (
	"container/heap"
	"fmt"
	"har"
	"net/http"
	"sync"
	"time"
)

type job struct {
	request  har.Request
	schedule time.Time
	index    int
}

type WorkerQueue []*job

func (jq WorkerQueue) Len() int { return len(jq) }

func (jq WorkerQueue) Less(i, j int) bool {
	return jq[i].schedule.Before(jq[j].schedule)
}

func (jq WorkerQueue) Swap(i, j int) {
	jq[i], jq[j] = jq[j], jq[i]
	jq[i].index = i
	jq[j].index = j
}

func (jq *WorkerQueue) Push(x interface{}) {
	n := len(*jq)
	job := x.(*job)
	job.index = n
	*jq = append(*jq, job)
}

func (jq *WorkerQueue) Pop() interface{} {
	old := *jq
	n := len(old)
	job := old[n-1]
	job.index = -1 // for safety
	*jq = old[0 : n-1]
	return job
}

type AsyncReplayer struct {
	Workers int
	Queue   WorkerQueue

	group  sync.WaitGroup
	closed bool
}

func NewAsyncReplayer(workers int) *AsyncReplayer {
	var ar AsyncReplayer
	ar.Workers = workers

	ar.Queue = make(WorkerQueue, 0)
	heap.Init(&ar.Queue)

	ar.closed = false

	for i := 0; i < workers; i++ {
		go ar.worker(i)
	}

	return &ar
}

func (ar *AsyncReplayer) Close() {
	ar.closed = true
}

func (ar *AsyncReplayer) Wait() {
	ar.group.Wait()
}

func (ar *AsyncReplayer) Replay(hardata *har.Har, opts *Options) {

	entries := hardata.Log.Entries

	if len(entries) == 0 {
		// Well, that was quick..
		return
	}

	start := time.Now()
	epoch := entries[0].Started

	ar.group.Add(len(entries))
	for _, entry := range entries {
		now := time.Now()
		delayFromStart := entry.Started.Sub(epoch)
		delay := delayFromStart - now.Sub(start)
		heap.Push(&ar.Queue, &job{
			request:  entry.Request,
			schedule: now.Add(delay),
		})
	}

	return
}

func (ar *AsyncReplayer) ReplayTimes(times int, hardata *har.Har, opts *Options) {
	for i := 0; i < times; i++ {
		ar.Replay(hardata, opts)
	}
}

func (ar *AsyncReplayer) worker(n int) {
	client := new(http.Client)

	var now time.Time
	var j *job
	for {
		if ar.closed {
			break
		}
		now = time.Now()
		j = nil
		if ar.Queue.Len() > 0 {
			j = heap.Pop(&ar.Queue).(*job)
		}
		if j == nil || j.schedule.After(now) {
			if j != nil {
				heap.Push(&ar.Queue, j)
			}
			time.Sleep(1 * time.Millisecond)
			continue
		}
		fmt.Println(n, "got", j.request.URL)
		Fire(client, &j.request)
		ar.group.Done()
	}
}
