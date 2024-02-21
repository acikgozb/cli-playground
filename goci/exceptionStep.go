package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type exceptionStep struct {
	step
}

// Since there are no new fields for this type, we can call newStep directly.
func newExceptionStep(name, exe, message, proj string, args []string) exceptionStep {
	s := exceptionStep{}
	s.step = newStep(name, exe, message, proj, args)
	return s
}

// Define a new version of execute() method that is different than what step{} has.
func (s exceptionStep) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &StepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	// If there is no error but output is not empty,
	// It means the error is captured on output.
	if out.Len() > 0 {
		return "", &StepErr{
			step:  s.name,
			msg:   fmt.Sprintf("invalid format: %s", out.String()),
			cause: nil,
		}
	}

	return s.message, nil
}
