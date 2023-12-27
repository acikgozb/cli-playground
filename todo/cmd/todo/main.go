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

	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil && !os.IsNotExist(err) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *listFlag:
		for _, item := range *l {
			if !item.Done {
				fmt.Println(item.Task)
			}
		}
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
