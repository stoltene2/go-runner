package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	timeout  <-chan time.Time
	interrupt chan os.Signal
	result   chan error
	tasks    []func() error
}

var ErrTimeout = errors.New("Timeout")
var ErrInterrupt = errors.New("Interrupt received")

func New(t time.Duration) *Runner {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	return &Runner{
		timeout:  time.After(t),
		result:   make(chan error),
		interrupt: interrupt,
	}
}

func (r *Runner) AddTasks(tasks ...func() error) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {
	go func() {
		r.result <- r.run()
	}()

	select {
	case <-r.timeout:
		return ErrTimeout
	case result := <- r.result:
		return result
	}
}

func (r *Runner) run() error {
	for _, t := range r.tasks {
		select {
		case <-r.interrupt:
			signal.Stop(r.interrupt)
			return ErrInterrupt
		default:
			t()
		}
	}

	return nil
}
