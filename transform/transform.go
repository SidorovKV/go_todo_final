package transform

import (
	"fmt"
	"go_todo_final/model"
	"go_todo_final/services/todolist"
	"strconv"
	"time"
)

func DtoToTask(task model.TaskDTO) (model.Task, error) {
	if task.Title == "" {
		return model.Task{}, fmt.Errorf("empty title")
	}

	today := makeDate(time.Now())
	date := today

	var err error

	rawDate := task.Date
	if rawDate != "" {
		date, err = time.Parse("20060102", rawDate)
		if err != nil {
			return model.Task{}, fmt.Errorf("invalid date format")
		}
	}

	if date.Before(today) {
		if task.Repeat == "" {
			date = today
		} else {
			date, err = todolist.NextDate(date, today, task.Repeat)
			if err != nil {
				return model.Task{}, fmt.Errorf("can't get next date: %w", err)
			}
		}
	}

	id := 0
	if task.ID != "" {
		id, _ = strconv.Atoi(task.ID)
	}

	return model.Task{
		ID:      uint64(id),
		Date:    date.Format("20060102"),
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}, nil
}

func TasksToDto(tasks []model.Task) []model.TaskDTO {
	dtos := make([]model.TaskDTO, 0, len(tasks))

	for _, task := range tasks {
		dtos = append(dtos, model.TaskDTO{
			ID:      strconv.Itoa(int(task.ID)),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	return dtos
}

func TaskToDto(task model.Task) model.TaskDTO {
	return model.TaskDTO{
		ID:      strconv.Itoa(int(task.ID)),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}

func makeDate(datetime time.Time) time.Time {
	y, m, d := datetime.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
