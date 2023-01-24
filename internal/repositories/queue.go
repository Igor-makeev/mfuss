package repositories

import (
	"context"
	"sync"
	"time"
)

type Queue struct {
	stream         chan string
	Buf            []string
	UpdateInterval time.Duration
	sync.Mutex
	sync.WaitGroup
}

func NewQueue() *Queue {
	q := &Queue{
		stream:         make(chan string, 15),
		Buf:            make([]string, 0, 50),
		UpdateInterval: 3 * time.Second}

	return q
}

func (q *Queue) Write(data []string) {
	q.WaitGroup.Add(1)
	go func() {
		for _, v := range data {
			q.stream <- v
		}
		q.Done()
	}()

}

func (q *Queue) CleanBuf() {
	q.Mutex.Lock()
	q.Buf = q.Buf[:0]
	q.Mutex.Unlock()
}
func (q *Queue) Close() {
	close(q.stream)
}

func (q *Queue) Listen(ctx context.Context, commitData func([]string) error, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		q.Wait()
		select {
		case <-q.stream:
			if len(q.Buf) == cap(q.Buf) {
				commitData(q.Buf)
			}
			q.CleanBuf()
			q.Buf = append(q.Buf, <-q.stream)
		case <-ticker.C:
			commitData(q.Buf)
			q.CleanBuf()
			continue
		case <-ctx.Done():
			commitData(q.Buf)
			q.CleanBuf()
			return
		}
	}
}
