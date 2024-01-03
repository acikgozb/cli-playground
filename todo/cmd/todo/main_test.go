package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	binName  = "todo"
	fileName = "todo.json"
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
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task1 := "test task number 1"
	task2 := "Hello from pipe"

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

		expected := fmt.Sprintf("   1: %s\n   2: %s\n", task1, task2)
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
}
