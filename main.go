package main

import "go-runner/runner"
import "time"
import "fmt"

func createTask(n int) func() error {
	return func() error {
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Println("Executing task: ", n)
		return nil
	}
}

func main() {

	// Create a runner that
	// *. will stop running when an interrupt is received

	r := runner.New(time.Duration(2) * time.Second)
	r.AddTasks(createTask(0), createTask(1), createTask(2), createTask(3), createTask(4))

	fmt.Println("Running")

	err := r.Start()

	if err != nil {
		fmt.Println("Caught error", err)
	} else {
		fmt.Println("All tasks exited successfully")
	}
}
