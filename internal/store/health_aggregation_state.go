package store

import "errors"

type HealthAggregationStream struct{}

func NewHealthAggregationStream() *HealthAggregationStream { return &HealthAggregationStream{} }

func (s *HealthAggregationStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 两条 channel 都必须在 goroutine 退出前关闭：
		// out 不关，消费端的 for range 永不返回——这正是“收到第一份就卡住”的死锁根因，
		// 错误路径早退时同样要关，才能让消费端终止并保留已收集的数据。
		defer close(out)
		defer close(errs)
		for index, item := range items {
			if index == failAt {
				// errs 缓冲为 1，写入立即成功；错误在 return 前入队，
				// 故消费端排空 out 后总能从此读到终止错误（与 close 顺序无关）。
				errs <- errors.New("stream interrupted")
				return
			}
			out <- item
		}
	}()
	return out, errs
}
