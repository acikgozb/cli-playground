package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	_, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not installed, the test is skipped")
	}

	testCases := []struct {
		name     string
		proj     string
		out      string
		expErr   error
		setupGit bool
	}{
		{
			name:     "success",
			proj:     "./testdata/tool/",
			out:      "Go build: SUCCESS\nGo test: SUCCESS\ngofmt: SUCCESS\ngit push: SUCCESS\n",
			expErr:   nil,
			setupGit: true,
		},
		{
			name:     "fail",
			proj:     "./testdata/toolErr/",
			out:      "",
			expErr:   &StepErr{step: "go build"},
			setupGit: false,
		},
		{
			name:     "failFormat",
			proj:     "./testdata/toolFmtErr/",
			out:      "",
			expErr:   &StepErr{step: "go fmt"},
			setupGit: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupGit {
				cleanup := setupGit(t, tc.proj)
				defer cleanup()
			}

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

// Helper function to create a mock Git environment to test the git push pipeline step
func setupGit(t *testing.T, proj string) func() {
	t.Helper()

	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}

	absProjPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	remoteGitUri := fmt.Sprintf("file://%s", tempDir)
	gitCmdList := []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, absProjPath, nil},
		{[]string{"remote", "add", "testOrigin", remoteGitUri}, absProjPath, nil},
		{[]string{"add", "."}, absProjPath, nil},
		{[]string{"commit", "-m", "test"}, absProjPath, []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
		}},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	// Return the cleanup function to use after tests.
	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(absProjPath, ".git"))
	}
}
