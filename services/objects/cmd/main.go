package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amir-mln/amdp-task/services/objects/cmd/http"
)

func broadcast[T any](count int, input <-chan T) []chan T {
	output := make([]chan T, count)
	for i := range count {
		output[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, outCh := range output {
				close(outCh)
			}
		}()

		for data := range input {
			for _, outCh := range output {
				outCh <- data
			}
		}
	}()

	return output
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sigChs := broadcast(2, sigCh)
	httpErrCh := make(chan error)

	go http.Run(sigChs[0], httpErrCh)
	select {
	case <-sigChs[len(sigChs)-1]:
		time.AfterFunc(6*time.Second, func() { os.Exit(1) })
	case err := <-httpErrCh:
		fmt.Println(err)
		os.Exit(1)
		// call other services to terminate
	}
}
