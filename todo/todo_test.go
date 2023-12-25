package todo_test

import (
	"github.com/acikgozb/cli-playground/todo"
	"os"
	"testing"
)

func TestList_Add(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("expected %s task name but got %s", taskName, l[0].Task)
	}
}

func TestList_Complete(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	if l[0].Done {
		t.Errorf("expected new task to not be completed")
	}

	err := l.Complete(1)
	if err != nil {
		t.Errorf("expected l.Complete to work but got err: %v", err)
	}

	if !l[0].Done {
		t.Errorf("expected task to be completed but got %v", l[0].Done)
	}
}

func TestList_Delete(t *testing.T) {
	l := todo.List{}

	tasks := []string{
		"New Task1",
		"New Task2",
		"New Task3",
	}

	for _, task := range tasks {
		l.Add(task)
	}

	err := l.Delete(2)
	if err != nil {
		t.Errorf("expected task to be removed but got err: %v", err)
	}

	if len(l) != 2 {
		t.Errorf("expected list length to be 2 but got %d", len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("expected %q after deleting the item, but got %q instead", tasks[2], l[1].Task)
	}
}

func TestList_Save_Get(t *testing.T) {
	list1 := todo.List{}
	list2 := todo.List{}

	taskName := "New Task"
	list1.Add(taskName)

	file, err := os.CreateTemp("", "test-list")
	if err != nil {
		t.Fatalf("unable to create file: %v", err)
	}

	defer os.Remove(file.Name())

	if err = list1.Save(file.Name()); err != nil {
		t.Errorf("unable to save the list to the file: %v", err)
	}

	if err = list2.Get(file.Name()); err != nil {
		t.Errorf("unable to get the list from the file: %v", err)
	}

	if list1[0].Task != list2[0].Task {
		t.Errorf("expected %q but got %q instead", list1[0].Task, list2[0].Task)
	}
}
