package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, todo)
}

func (l *List) Complete(itemNumber int) error {
	if itemNumber > len(*l) || itemNumber <= 0 {
		return fmt.Errorf("item %d does not exist", itemNumber)
	}

	list := *l

	list[itemNumber-1].Done = true
	list[itemNumber-1].CompletedAt = time.Now()

	return nil
}

func (l *List) Delete(itemNumber int) error {
	if itemNumber > len(*l) || itemNumber <= 0 {
		return fmt.Errorf("item %d does not exist", itemNumber)
	}

	list := *l
	list = append(list[:itemNumber-1], list[itemNumber:]...)
	*l = list

	return nil
}

func (l *List) Save(fileName string) error {
	listJSON, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, listJSON, 0644)
}

func (l *List) Get(fileName string) error {
	listJSON, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if len(listJSON) == 0 {
		return nil
	}

	return json.Unmarshal(listJSON, l)
}

func (l *List) String() string {
	formatted := ""

	for index, item := range *l {
		prefix := "   "
		if item.Done {
			prefix = "X  "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, index+1, item.Task)
	}

	return formatted
}
