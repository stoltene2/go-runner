package runner

import (
	"fmt"
)

type Runner struct {

	tasks []func() error
}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) AddTasks(tasks ...func() error) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {
	fmt.Println("Running")

	for _, t := range r.tasks {
		t()
	}

	return nil
}
