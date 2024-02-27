package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/acikgozb/cli-playground/pscan/scan"
)

// Since host actions depend on a host file to work on
// We create a helper method to generate a pre-filled dummy host file for us.
func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	temp, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}

	// We directly close the file since actions just need the file name.
	temp.Close()

	if initList {
		hl := &scan.HostsList{}

		for _, host := range hosts {
			hl.Add(host)
		}

		if err := hl.Save(temp.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// Return the temp file name and cleanup func.
	return temp.Name(), func() {
		os.Remove(temp.Name())
	}
}

func TestHostActions(t *testing.T) {
	// Define hosts for action test.
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOutput: "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			args:           hosts,
			expectedOutput: "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOutput: "Removed the host: host1\nRemoved the host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup action test.
			tempFileName, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			// Define output to capture
			var out bytes.Buffer

			// Execute action and capture the output
			if err := tc.actionFunction(&out, tempFileName, tc.args); err != nil {
				t.Fatalf("expected no error but got %q\n", err)
			}

			// Verify the output.
			if out.String() != tc.expectedOutput {
				t.Errorf("expected %q as a result but got %q\n", tc.expectedOutput, out.String())
			}
		})
	}
}

// Integration test
// The goal is to execute all commands in sequence, simulating what a user would do.
// Flow: Add 3 hosts, list them and delete one host from the list.
func TestIntegration(t *testing.T) {
	// Define the hosts for the test.
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Setup the integration test.
	// We do not want to initialize the list, CLI should be able to do it by itself.
	tempFileName, cleanup := setup(t, hosts, false)
	defer cleanup()

	hostToDelete := "host2"
	hostsAfterDeletion := []string{
		"host1",
		"host3",
	}

	// Capture the output
	var out bytes.Buffer

	// Define the expected output by combining all outputs we expect from the test.
	expectedOutput := ""

	// Output after inserting hosts.
	for _, host := range hosts {
		expectedOutput += fmt.Sprintf("Added host: %s\n", host)
	}

	// Output after listing hosts.
	expectedOutput += strings.Join(hosts, "\n")
	expectedOutput += fmt.Sprintln()

	// Output after deleting a host.
	expectedOutput += fmt.Sprintf("Removed the host: %s\n", hostToDelete)
	expectedOutput += strings.Join(hostsAfterDeletion, "\n")
	expectedOutput += fmt.Sprintln()

	// Execute all operations in defined sequence add > list > delete > list.

	// Add hosts to the list.
	if err := addAction(&out, tempFileName, hosts); err != nil {
		t.Fatalf("expected no error from addAction but got %q instead", err)
	}

	// List hosts.
	if err := listAction(&out, tempFileName, nil); err != nil {
		t.Fatalf("expected no error from listAction but got %q instead", err)
	}

	// Delete a host from the list.
	if err := deleteAction(&out, tempFileName, []string{hostToDelete}); err != nil {
		t.Fatalf("expected no error from deleteAction but got %q instead", err)
	}

	// List remaining hosts.
	if err := listAction(&out, tempFileName, nil); err != nil {
		t.Fatalf("expected no error from listAction after deletion but got %q instead", err)
	}

	// Test the output
	if out.String() != expectedOutput {
		t.Errorf("expected output to be %q, but got %q instead", expectedOutput, out.String())
	}
}
