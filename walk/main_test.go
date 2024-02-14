package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name:     "NoFilter",
			root:     "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name:     "FilterExtensionMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n",
		},
		{
			name:     "FilterExtensionSizeMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n",
		},
		{
			name:     "FilterExtensionSizeNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: "",
		},
		{
			name:     "FilterExtensionNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(testCase.root, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if testCase.expected != res {
				t.Errorf("Expected %q, got %q instead", testCase.expected, res)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{
			name:        "DeleteExtensionNoMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			expected:    "",
		},
		{
			name:        "DeleteExtensionMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			expected:    "",
		},
		{
			name:        "DeleteExtensionMixed",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     5,
			nNoDelete:   5,
			expected:    "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			var logBuffer bytes.Buffer

			testCase.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.cfg.ext:     testCase.nDelete,
				testCase.extNoDelete: testCase.nNoDelete,
			})

			defer cleanup()

			if err := run(tempDir, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}

			result := buffer.String()

			if testCase.expected != result {
				t.Errorf("Expected %q, got %q instead\n", testCase.expected, result)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != testCase.nNoDelete {
				t.Errorf(
					"Expected %d files left, but got %d instead\n",
					testCase.nNoDelete,
					len(filesLeft),
				)
			}

			expectedLogLines := testCase.nDelete + 1 // +1 corresponds to the empty final line in log file.
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expectedLogLines {
				t.Errorf("Expected %d log lines, got %d instead", expectedLogLines, len(lines))
			}
		})
	}
}

func TestRunArchive(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{
			name:         "ArchiveExtensionNoMatch",
			cfg:          config{ext: ".log"},
			extNoArchive: ".gz",
			nArchive:     0,
			nNoArchive:   10,
		},
		{
			name:         "ArchiveExtensionMatch",
			cfg:          config{ext: ".log"},
			extNoArchive: "",
			nArchive:     10,
			nNoArchive:   0,
		},
		{
			name:         "ArchiveExtensionMixed",
			cfg:          config{ext: ".log"},
			extNoArchive: ".gz",
			nArchive:     5,
			nNoArchive:   5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Buffer for RunArchive output
			var buffer bytes.Buffer

			// Create temp dirs to RunArchive test
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:      tc.nArchive,
				tc.extNoArchive: tc.nNoArchive,
			})

			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tc.cfg.archive = archiveDir

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			// Validate the output
			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
			expectedFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}

			expected := strings.Join(expectedFiles, "\n")
			// Remove last line from the buffer before comparing the outputs
			result := strings.TrimSpace(buffer.String())

			if expected != result {
				t.Errorf("Expected %q, got %q instead", expected, result)
			}

			// Validate the number of archived files
			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.nArchive {
				t.Errorf(
					"Expected %d archived files, got %d instead\n",
					tc.nArchive,
					len(filesArchived),
				)
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fileName := fmt.Sprintf("file%d%s", j, k)
			filePath := filepath.Join(tempDir, fileName)
			if err := os.WriteFile(filePath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}
