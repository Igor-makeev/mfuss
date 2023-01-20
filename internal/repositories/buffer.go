package repositories

import (
	"sync"
	"time"
)

type Buffer struct {
	stream         chan string
	buf            []string
	updateInterval time.Duration
	sync.Mutex
}

func NewBuffer() *Buffer {
	bf := &Buffer{
		stream:         make(chan string, 10),
		buf:            make([]string, 0, 10),
		updateInterval: 5 * time.Second}

	return bf
}

func (b *Buffer) Write(str string) {

	b.stream <- str

}

func (b *Buffer) Read() string {

	return <-b.stream
}
func (b *Buffer) CleanBuf() {
	b.Mutex.Lock()
	b.buf = b.buf[:0]
	b.Mutex.Unlock()
}
func (b *Buffer) Close() {
	close(b.stream)
}

// func (b *Buffer) startWork(ctx context.Context) {
// 	for {

// 		select {
// 		case <-time.After(b.updateInterval):
// 			// в базу и чистим
// 		case <-b.stream:
// 			// в базу и чистим
// 		case <-ctx.Done():
// 		}

// 	}
// }

// func fanIn(bf *Buffer, inputChs ...chan string) chan string {

// 	go func() {
// 		wg := &sync.WaitGroup{}

// 		for _, inputCh := range inputChs {
// 			wg.Add(1)

// 			go func(inputCh chan string) {
// 				defer wg.Done()
// 				for item := range inputCh {
// 					bf.stream <- item
// 				}
// 			}(inputCh)
// 		}

// 		wg.Wait()
// 		close(bf.stream)
// 	}()

// 	return bf.stream
// }
