package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/savioxavier/termlink"
	"google.golang.org/api/calendar/v3"
	"os"
	"strings"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	Hide        bool
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Hide:        false,
	}
	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	completedTime := time.Now()
	list[index-1].CompletedAt = &completedTime
	list[index-1].Done = true
	return nil
}

func (t *Todos) Uncomplete(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	list[index-1].CompletedAt = nil
	list[index-1].Done = false
	return nil
}

func (t *Todos) Delete(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	*t = append(list[:index-1], list[index:]...)
	return nil
}

func (t *Todos) MoveStart(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	itemToMove := list[index-1]
	list = append(list[:index-1], list[index:]...)
	newList := make([]item, 0)
	newList = append(newList, itemToMove)
	for _, item := range list {
		newList = append(newList, item)
	}
	*t = newList
	return nil
}

func (t *Todos) MoveEnd(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	itemToMove := list[index-1]
	newList := make([]item, 0)
	list = append(list[:index-1], list[index:]...)
	for _, item := range list {
		newList = append(newList, item)
	}
	newList = append(newList, itemToMove)
	*t = newList
	return nil
}

func (t *Todos) Hide(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	list[index-1].Hide = true
	return nil
}

func (t *Todos) HideAll() error {
	list := *t
	for i, item := range list {
		if item.Done == true && item.CompletedAt != nil {
			list[i].Hide = true
		}
	}
	return nil
}

func (t *Todos) Unhide(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("invalid index")
	}
	list[index-1].Hide = false
	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) CheckForDoneItemsYesterday(filename string) {
	_ = t.Load(filename)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	list := *t
	for i, item := range list {
		if item.CompletedAt != nil && !item.Hide {
			if item.CompletedAt.Before(today) {
				list[i].Hide = true
			}
		}
	}
	_ = t.Store(filename)
}

func (t *Todos) Print(all bool, events *calendar.Events) {
	tableCalendar := simpletable.New()
	tableCalendar.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Event"},
			{Align: simpletable.AlignCenter, Text: "Start"},
			{Align: simpletable.AlignCenter, Text: "End"},
			{Align: simpletable.AlignCenter, Text: "Link"},
		},
	}

	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found in your google calendar.")
	} else {
		fmt.Println("Upcoming events from Google Calendar:")
	}

	var cellsCalendar [][]*simpletable.Cell
	for i, item := range events.Items {
		i++
		event := green(item.Summary)
		startDate := item.Start.DateTime
		if startDate == "" {
			startDate = item.Start.Date
		} else {
			startDate = strings.ReplaceAll(startDate, "T", " ")
			startDate = startDate[:len(startDate)-9]
		}
		endDate := item.End.DateTime
		if endDate == "" {
			endDate = "(all day)"
		} else {
			endDate = strings.ReplaceAll(endDate, "T", " ")
			endDate = endDate[:len(endDate)-9]
		}
		cellsCalendar = append(cellsCalendar, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i)},
			{Text: event},
			{Text: startDate},
			{Text: endDate},
			{Text: termlink.ColorLink("@", item.HtmlLink, "italic green", true)},
		})
	}

	tableCalendar.Body = &simpletable.Body{Cells: cellsCalendar}
	tableCalendar.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: ""},
		},
	}
	tableCalendar.SetStyle(simpletable.StyleCompactLite)
	tableCalendar.Println()

	fmt.Println()
	fmt.Println("Your tasks:")
	hiddenItems := 0
	completedItems := 0
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "Created at"},
			{Align: simpletable.AlignCenter, Text: "Completed at"},
		},
	}
	var cells [][]*simpletable.Cell
	for i, item := range *t {
		i++
		task := blue(item.Task)
		if item.Done {
			task = gray(item.Task)
		}
		if item.Hide {
			hiddenItems++
			if !all {
				continue
			}
		}
		if item.Done {
			completedItems++
		}
		var completedAt string
		if item.CompletedAt != nil {
			completedAt = item.CompletedAt.Format("2006-01-02 15:04")
		} else {
			completedAt = ""
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i)},
			{Text: task},
			{Text: fmt.Sprintf("%t", item.Done)},
			{Text: item.CreatedAt.Format("2006-01-02 15:04")},
			{Text: completedAt},
		})
	}
	allItems := len(*t)

	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: fmt.Sprintf("Total: %v, Completed: %v, Hidden: %v", allItems, completedItems, hiddenItems)},
		},
	}
	table.SetStyle(simpletable.StyleCompactLite)
	table.Println()
}
