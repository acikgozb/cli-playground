package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	pipeline := make([]executer, 3)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep("go test", "go", "Go test: SUCCESS", proj, []string{"test", "-v"})
	pipeline[2] = newExceptionStep("go fmt", "gofmt", "gofmt: SUCCESS", proj, []string{"-l", "."})

	for _, step := range pipeline {
		msg, err := step.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
