package main

import (
	"context"
	"os/exec"
	"time"
)

var command = exec.CommandContext

type timeoutStep struct {
	step
	timeout time.Duration
}

func newTimeoutStep(
	name, exe, message, proj string,
	args []string,
	timeout time.Duration,
) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, exe, message, proj, args)
	s.timeout = timeout
	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}

	return s
}

func (s timeoutStep) execute() (string, error) {
	// Create a context to apply timeout in case goci hangs
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := command(ctx, s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &StepErr{
				step:  s.name,
				msg:   "timeout, command failed",
				cause: context.DeadlineExceeded,
			}
		}

		return "", &StepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
