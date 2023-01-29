package service

import (
	"context"
	"sync"
	"time"
)

var bufCap = 200

type Queue struct {
	stream chan string

	Buf            []string
	UpdateInterval time.Duration
	mutex          sync.Mutex
}

func NewQueue() *Queue {
	q := &Queue{
		stream:         make(chan string, 50),
		Buf:            make([]string, 0, bufCap),
		UpdateInterval: 1 * time.Second,
	}

	return q
}

func (q *Queue) Run(ctx context.Context, commitData func(context.Context, []string) error) {
	go q.listen(ctx, commitData, q.UpdateInterval)
}

func (q *Queue) Write(data []string) {

	go func([]string) {
		for _, v := range data {
			q.stream <- v
		}

	}(data)

}

func (q *Queue) cleanBuf() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.Buf = q.Buf[:0]

}
func (q *Queue) Close() {
	close(q.stream)
}

func (q *Queue) commitData(ctx context.Context, commitDataFunc func(context.Context, []string) error) {
	commitDataFunc(ctx, q.Buf)
	q.cleanBuf()

}

func (q *Queue) listen(ctx context.Context, commitDataFunc func(context.Context, []string) error, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {

		select {
		case id := <-q.stream:
			if len(q.Buf) == bufCap {
				q.commitData(ctx, commitDataFunc)
			}

			q.Buf = append(q.Buf, id)
			continue
		case <-ticker.C:

			q.commitData(ctx, commitDataFunc)
			continue
		case <-ctx.Done():

			q.commitData(ctx, commitDataFunc)
			return
		}
	}
}
