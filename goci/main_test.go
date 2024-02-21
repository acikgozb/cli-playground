package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		out    string
		expErr error
	}{
		{
			name:   "success",
			proj:   "./testdata/tool/",
			out:    "Go build: SUCCESS\nGo test: SUCCESS\n",
			expErr: nil,
		},
		{
			name:   "fail",
			proj:   "./testdata/toolErr/",
			out:    "",
			expErr: &StepErr{step: "go build"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(tc.proj, &out)

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Excepted error: %q. Got nil instead", tc.expErr)
					return
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Excepted error: %q. Got %q", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != tc.out {
				t.Errorf("Excepted output: %q. Got %q instead", tc.out, out.String())
			}
		})
	}
}
