package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	pipeline := make([]executer, 4)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep("go test", "go", "Go test: SUCCESS", proj, []string{"test", "-v"})
	pipeline[2] = newExceptionStep("go fmt", "gofmt", "gofmt: SUCCESS", proj, []string{"-l", "."})
	pipeline[3] = newTimeoutStep(
		"git push",
		"git",
		"git push: SUCCESS",
		proj,
		[]string{"push", "origin", "master"},
		10*time.Second,
	)

	// Create a buffered channel to capture os signal changes and act accordingly.
	signalChannel := make(chan os.Signal, 1)
	errCh := make(chan error)
	done := make(chan struct{})

	// Pass in the interested signals to the channel
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, step := range pipeline {
			msg, err := step.execute()
			if err != nil {
				errCh <- err
				return
			}

			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errCh <- err
				return
			}
		}
		close(done)
	}()

	// Handle each output from channels accordingly.
	for {
		select {
		case rec := <-signalChannel:
			signal.Stop(signalChannel)
			return fmt.Errorf("%s: Exiting: %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}
