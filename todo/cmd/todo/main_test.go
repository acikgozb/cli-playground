package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	binName         = "todo"
	defaultFileName = "todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building the tool...")

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot build the binary %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)

	var fileName string
	if os.Getenv("TODO_FILENAME") == "" {
		fileName = defaultFileName
	}

	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task1 := "mockAddNewTask item"
	task2 := "mockAddNewTaskFromSTDIN item"
	task3 := "mockShowNonCompletedTasks item"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task1)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		listCmd := exec.Command(cmdPath, "-list")
		out, err := listCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		completedMark := "X"
		if !strings.Contains(string(out), completedMark) {
			t.Errorf("Expected %s mark on task but got: %s", completedMark, string(out))
		}
	})

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Errorf("Error occured when creating a stdin pipe: %s", err.Error())
		}

		_, err = io.WriteString(cmdStdIn, task2)
		if err != nil {
			t.Errorf("Error occured while writing to a pipe: %s", err.Error())
		}

		if err = cmdStdIn.Close(); err != nil {
			t.Fatal(err)
		}

		if err = cmd.Run(); err != nil {
			t.Errorf("Error occured while executing the -add command: %s", err.Error())
		}
	})

	t.Run("ListExistingTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("X  1: %s\n   2: %s\n", task1, task2)
		if expected != string(out) {
			t.Errorf("Expected %s from cli but got %s", expected, string(out))
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		deleteCmd := exec.Command(cmdPath, "-delete", "1")
		if err := deleteCmd.Run(); err != nil {
			t.Fatal(err)
		}

		listCmd := exec.Command(cmdPath, "-list")
		out, err := listCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("   1: %s\n", task2)
		if expected != string(out) {
			t.Errorf("Expected %s from cli but got %s", expected, string(out))
		}
	})

	t.Run("ShowListAsVerbose", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		tokenStartIndex := strings.Index(string(out), "[{")
		tokenEndIndex := strings.Index(string(out), "}]")

		if tokenStartIndex != 0 {
			t.Errorf("Expected %d as starting index, but got %d", 0, tokenStartIndex)
		}

		if tokenEndIndex == len(string(out))-1 {
			t.Errorf("Expected }] at the end to contain in the output, but got %s", string(out))
		}
	})

	t.Run("ShowNonCompletedTasks", func(t *testing.T) {
		addCmd := exec.Command(cmdPath, "-add", task3)
		if err := addCmd.Run(); err != nil {
			t.Fatal(err)
		}

		completeCmd := exec.Command(cmdPath, "-complete", "2")
		if err := completeCmd.Run(); err != nil {
			t.Fatal(err)
		}

		ncCmd := exec.Command(cmdPath, "-nc")
		out, err := ncCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("   1: %s\n", task2)
		if expected != string(out) {
			t.Errorf("Expected to only see the not completed task, but got %s", string(out))
		}
	})
}
