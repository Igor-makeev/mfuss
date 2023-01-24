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
}

func NewQueue() *Queue {
	q := &Queue{
		stream:         make(chan string, 50),
		Buf:            make([]string, 0, 200),
		UpdateInterval: 3 * time.Second}

	return q
}

func (q *Queue) Write(data []string) {

	go func([]string) {
		for _, v := range data {
			q.stream <- v
		}

	}(data)

}

func (q *Queue) CleanBuf() {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	q.Buf = q.Buf[:0]

}
func (q *Queue) Close() {
	close(q.stream)
}

func (q *Queue) Listen(ctx context.Context, commitData func([]string) error, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {

		select {
		case id := <-q.stream:
			if len(q.Buf) == cap(q.Buf) {
				commitData(q.Buf)
				q.CleanBuf()
			}

			q.Buf = append(q.Buf, id)
			continue
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
