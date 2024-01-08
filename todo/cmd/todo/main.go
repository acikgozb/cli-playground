package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/acikgozb/cli-playground/todo"
	"io"
	"os"
)

var todoFileName = "todo.json"

func main() {
	addFlag := flag.Bool("add", false, "Adds given task to the todo list")
	listFlag := flag.Bool("list", false, "Lists all todo items that are not completed")
	completeFlag := flag.Int("complete", 0, "Marks the item with given itemNumber as completed")
	deleteFlag := flag.Int("delete", 0, "Deletes the item with given itemNumber from the todo list")
	verboseFlag := flag.Bool("verbose", false, "Enables verbose output, shows the item in JSON format")
	showNonCompletedFlag := flag.Bool("nc", false, "Shows tasks which are not completed")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Todo tool. Developed by acikgozb.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "The CLI is NOT production ready, keep this in mind while using.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "To add a new task, simply enter your task with the -add option:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "todo -add \"My new task\"\n")
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
	case *verboseFlag:
		jsonOut, err := verboseOut(l)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(jsonOut)
	case *showNonCompletedFlag:
		ncTasks := l.NotCompletedTasks()
		fmt.Print(ncTasks)
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
		tasks, err := getTasks(os.Stdin, flag.Args()...)
		fmt.Println(tasks)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Add(tasks)

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

func getTasks(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return args, nil
	}

	var tasks []string

	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		task := s.Text()
		if len(task) == 0 {
			return nil, fmt.Errorf("task cannot be blank")
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func verboseOut(list *todo.List) (string, error) {
	byteOut, err := json.Marshal(list)
	if err != nil {
		return "", err
	}

	return string(byteOut), nil
}
