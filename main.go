package main

import "go-runner/runner"
import "time"
import "fmt"

func createTask(n int) func() error {
	return func () error {
		time.Sleep(time.Second)
		fmt.Println("Executing task: ", n)
		return nil
	}
}

func main() {

	// Create a runner that
	// *. will stop running when all tasks complete
	// *. will stop running when a timeout happens
	// *. will stop running when an interrupt is received
	// *. will stop running if an error is thrown by a task


	r := runner.New()
	r.AddTasks(createTask(0), createTask(1), createTask(2))

	err := r.Start()

	if err != nil {
		fmt.Println("Caught error", err)
	} else {
		fmt.Println("All tasks exited successfully")
	}
}
