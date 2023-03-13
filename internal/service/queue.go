package service

import (
	"context"
	"mfuss/internal/repositories"
	"sync"
	"time"
)

// вместимость буфера
var bufCap = 200

// интерфейс комитера
type Commiter interface {
	MarkAsDeleted(ctx context.Context, arr []string) error
}

// структура очереди
type Queue struct {
	stream         chan string
	store          Commiter
	buf            []string
	updateInterval time.Duration
	mutex          sync.Mutex
}

// конструктор очереди
func NewQueue(rep *repositories.Repository) *Queue {
	q := &Queue{
		stream:         make(chan string, 50),
		store:          rep.URLStorager,
		buf:            make([]string, 0, bufCap),
		updateInterval: 1 * time.Second,
	}

	return q
}

// запустить очередь
func (q *Queue) Run(ctx context.Context) {
	go q.listen(ctx, q.updateInterval)
}

// записать в очередь
func (q *Queue) Write(data []string) {

	go func([]string) {
		for _, v := range data {
			q.stream <- v
		}

	}(data)

}

// очистить буфер
func (q *Queue) cleanBuf() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.buf = q.buf[:0]

}

// Закрыть поток
func (q *Queue) Close() {
	close(q.stream)
}

// закоммитить данные
func (q *Queue) commitData(ctx context.Context, arr []string) {
	q.store.MarkAsDeleted(ctx, arr)
	q.cleanBuf()

}

// функция которая слушает канал очереди
func (q *Queue) listen(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {

		select {
		case id := <-q.stream:
			if len(q.buf) == bufCap {
				q.commitData(ctx, q.buf)
			}

			q.buf = append(q.buf, id)
			continue
		case <-ticker.C:

			q.commitData(ctx, q.buf)
			continue
		case <-ctx.Done():

			q.commitData(ctx, q.buf)
			return
		}
	}
}
