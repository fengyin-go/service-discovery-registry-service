package store

import "errors"

type ProbeBatchStream struct{}

func NewProbeBatchStream() *ProbeBatchStream { return &ProbeBatchStream{} }

func (s *ProbeBatchStream) Start(items []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 无论正常结束还是遇错中断，都要关闭结果流，
		// 否则消费方的 range 会永久阻塞，导致批次挂起。
		defer close(out)
		defer close(errs)
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
