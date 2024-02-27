package scan_test

import (
	"errors"
	"os"
	"testing"

	"github.com/acikgozb/cli-playground/pscan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name           string
		host           string
		expectedLength int
		expectedError  error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			// Given
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			// When
			err := hl.Add(tc.host)

			// Then
			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error, got nil instead.\n")
				}

				if !errors.Is(err, tc.expectedError) {
					t.Errorf("expected error %q, got %q instead", tc.expectedError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error but got %q instead\n", err)
			}

			if len(hl.Hosts) != tc.expectedLength {
				t.Errorf(
					"expected list length %d, got %d instead",
					tc.expectedLength,
					len(hl.Hosts),
				)
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf(
					"expected host name %q as index 1, but got %q instead",
					tc.host,
					hl.Hosts[1],
				)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name           string
		host           string
		expectedLength int
		expectedError  error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			hl := &scan.HostsList{}

			for _, dummyHost := range []string{"host1", "host2"} {
				if err := hl.Add(dummyHost); err != nil {
					t.Fatal(err)
				}
			}

			// When
			err := hl.Remove(tc.host)

			// Then
			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("expected error but got nil instead\n")
				}

				if !errors.Is(err, tc.expectedError) {
					t.Errorf("expected error %q but got %q instead\n", tc.expectedError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error but got %q instead\n", err)
			}

			if len(hl.Hosts) != tc.expectedLength {
				t.Errorf(
					"expected list length as %d but got %d instead\n",
					tc.expectedLength,
					len(hl.Hosts),
				)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	// Given
	hl1 := &scan.HostsList{}
	hl2 := &scan.HostsList{}

	hostName := "host1"
	if err := hl1.Add(hostName); err != nil {
		t.Fatalf("expected no error while initializing lists in TestSaveLoad but got %q", err)
	}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("expected no error while creating a temp file for TestSaveLoad but got %q", err)
	}

	defer os.Remove(tempFile.Name())

	// When
	if err := hl1.Save(tempFile.Name()); err != nil {
		// Then
		t.Fatalf("expected no error while saving hl1 to tempFile but got %q instead", err)
	}

	// When
	if err := hl2.Load(tempFile.Name()); err != nil {
		// Then
		t.Fatalf("expected no error while loading hl2 with tempFile but got %q instead", err)
	}

	// Then
	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf(
			"expected two lists to have the same hosts but got %q in hl1, %q in hl2",
			hl1.Hosts[0],
			hl2.Hosts[0],
		)
	}
}

func TestLoadNotExist(t *testing.T) {
	hl := &scan.HostsList{}
	fileName := "file-which-does-not-exist"

	err := hl.Load(fileName)
	if err != nil {
		t.Errorf("expected load to not return an error if file does not exist, but got %q", err)
	}
}
