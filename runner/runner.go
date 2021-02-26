package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	timeout  <-chan time.Time
	interrupt <-chan os.Signal
	result   chan error
	complete chan bool
	tasks    []func() error
}

var ErrTimeout = errors.New("Timeout")
var ErrInterrupt = errors.New("Interrupt received")

func New(t time.Duration) *Runner {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	return &Runner{
		timeout:  time.After(t),
		result:   make(chan error, 1),
		interrupt: interrupt,
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
			case <-r.interrupt:
				r.result <- ErrInterrupt
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
