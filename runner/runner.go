package runner

import (
	"errors"
	"time"
)

type Runner struct {
	timeout  <-chan time.Time
	result   chan error
	complete chan bool
	tasks    []func() error
}

var ErrTimeout = errors.New("Timeout")

func New(t time.Duration) *Runner {
	return &Runner{
		timeout:  time.After(t),
		result:   make(chan error, 1),
		complete: make(chan bool, 1),
	}
}

func (r *Runner) AddTasks(tasks ...func() error) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {
	go func() {
		for _, t := range r.tasks {
			select {
			case <-r.timeout:
				r.result <- ErrTimeout
				break
			default:
				t()
			}
		}
	}()

	select {
	case <-r.timeout:
		return ErrTimeout
	case result := <- r.result:
		return result
	}
}
