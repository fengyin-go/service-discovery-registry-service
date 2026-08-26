package store

import "errors"

type InstanceDrainStream struct{}

func NewInstanceDrainStream() *InstanceDrainStream { return &InstanceDrainStream{} }

func (s *InstanceDrainStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 结束通道：无论成功还是错误分支，都关闭 out，使消费端的 range 退出而非永久等待。
		defer close(out)
		defer close(errs)
		for index, item := range items {
			if index == failAt {
				errs <- errors.New("stream interrupted")
				// 错误分支直接返回，由 defer 关闭 out；已成功发送的项仍保留在 out 中供消费端收集。
				return
			}
			out <- item
		}
	}()
	return out, errs
}
