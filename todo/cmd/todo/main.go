package main

import (
	"flag"
	"fmt"
	"github.com/acikgozb/cli-playground/todo"
	"os"
)

const todoFileName = "todo.json"

func main() {
	taskFlag := flag.String("task", "", "Adds given task to the todo list")
	listFlag := flag.Bool("list", false, "Lists all todo items that are not completed")
	completeFlag := flag.Int("complete", 0, "Marks given itemNumber as completed")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Todo tool. Developed by acikgozb.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "The CLI is NOT production ready, keep this in mind while using.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

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
	case *taskFlag != "":
		l.Add(*taskFlag)
		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		_, _ = fmt.Fprintln(os.Stderr, "An invalid operation is passed.")
		os.Exit(1)
	}
}
