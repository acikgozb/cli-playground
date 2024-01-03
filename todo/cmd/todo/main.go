package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/acikgozb/cli-playground/todo"
	"io"
	"os"
	"strings"
)

var todoFileName = "todo.json"

func main() {
	addFlag := flag.Bool("add", false, "Adds given task to the todo list")
	listFlag := flag.Bool("list", false, "Lists all todo items that are not completed")
	completeFlag := flag.Int("complete", 0, "Marks the item with given itemNumber as completed")
	deleteFlag := flag.Int("delete", 0, "Deletes the item with given itemNumber from the todo list")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Todo tool. Developed by acikgozb.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "The CLI is NOT production ready, keep this in mind while using.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil && !os.IsNotExist(err) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *listFlag:
		fmt.Print(l)
	case *completeFlag > 0:
		if err := l.Complete(*completeFlag); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *addFlag:
		task, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Add(task)

		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *deleteFlag > 0:
		if err := l.Delete(*deleteFlag); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		_, _ = fmt.Fprintln(os.Stderr, "An invalid operation is passed.")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
