package service

import (
	"context"
	"mfuss/internal/repositories"
	"sync"
	"time"
)

var bufCap = 200

type Commiter interface {
	MarkAsDeleted(ctx context.Context, arr []string) error
}

type Queue struct {
	stream         chan string
	store          Commiter
	Buf            []string
	UpdateInterval time.Duration
	mutex          sync.Mutex
}

func NewQueue(rep *repositories.Repository) *Queue {
	q := &Queue{
		stream:         make(chan string, 50),
		store:          rep.URLStorager,
		Buf:            make([]string, 0, bufCap),
		UpdateInterval: 1 * time.Second,
	}

	return q
}

func (q *Queue) Run(ctx context.Context) {
	go q.listen(ctx, q.UpdateInterval)
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

func (q *Queue) commitData(ctx context.Context, arr []string) {
	q.store.MarkAsDeleted(ctx, arr)
	q.cleanBuf()

}

func (q *Queue) listen(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {

		select {
		case id := <-q.stream:
			if len(q.Buf) == bufCap {
				q.commitData(ctx, q.Buf)
			}

			q.Buf = append(q.Buf, id)
			continue
		case <-ticker.C:

			q.commitData(ctx, q.Buf)
			continue
		case <-ctx.Done():

			q.commitData(ctx, q.Buf)
			return
		}
	}
}
