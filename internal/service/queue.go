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

func (q *Queue) Run(ctx context.Context, commitData func([]string) error) {
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

func (q *Queue) commitData(commitDataFunc func([]string) error) {
	commitDataFunc(q.Buf)
	q.cleanBuf()

}

func (q *Queue) listen(ctx context.Context, commitDataFunc func([]string) error, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {

		select {
		case id := <-q.stream:
			if len(q.Buf) == bufCap {
				q.commitData(commitDataFunc)
			}

			q.Buf = append(q.Buf, id)
			continue
		case <-ticker.C:
			commitDataFunc(q.Buf)
			q.commitData(commitDataFunc)
			continue
		case <-ctx.Done():
			commitDataFunc(q.Buf)
			q.commitData(commitDataFunc)
			return
		}
	}
}
