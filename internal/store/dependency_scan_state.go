package store

import "errors"

type DependencyScanStream struct{}

func NewDependencyScanStream() *DependencyScanStream { return &DependencyScanStream{} }

func (s *DependencyScanStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 先关闭数据通道，让消费者（Collect 的 range）及时结束；
		// 再关闭错误通道。两条路径（失败提前 return 或正常跑完）都关闭 out，
		// 否则消费者会因 out 永不关闭而卡在 range，无法返回部分结果与错误。
		defer close(errs)
		defer close(out)
		for index, item := range items {
			if index == failAt {
				errs <- errors.New("stream interrupted")
				return
			}
			out <- item
		}
	}()
	return out, errs
}
