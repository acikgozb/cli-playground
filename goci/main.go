package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

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

	// The first step of the pipeline is to build the project using go build.
	// The errors package is used to prevent go build from creating an executable file.
	// Go build does not create an executable if it is used for building multiple packages.
	args := []string{"build", ".", "errors"}

	buildCmd := exec.Command("go", args...)
	buildCmd.Dir = proj

	if err := buildCmd.Run(); err != nil {
		return &StepErr{step: "go build", msg: "go build failed", cause: err}
	}

	_, err := fmt.Fprintln(out, "Go build: SUCCESS")
	return err
}
