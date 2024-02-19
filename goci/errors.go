package main

import (
	"errors"
	"fmt"
)

var ErrValidation = errors.New("validation failed")

// Define a custom error to centralize the errors in each step of the CI pipeline
type StepErr struct {
	step  string
	msg   string
	cause error
}

// This method is added to the struct to satisfy the Error interface
func (s *StepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

// This method is added to see whether the target error matches a StepErr
func (s *StepErr) Is(target error) bool {
	t, ok := target.(*StepErr)
	if !ok {
		return false
	}

	return t.step == s.step
}

// errors.Is might try to Unwrap the underlying error, therefore implement Unwrap as well
func (s *StepErr) Unwrap() error {
	return s.cause
}
