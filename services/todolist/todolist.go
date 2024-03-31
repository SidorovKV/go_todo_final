package todolist

import (
	"errors"
	"fmt"
	"go_todo_final/model"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	daysLimit  = 400
	daysParam  = "d"
	weekParam  = "w"
	monthParam = "m"
	yearParam  = "y"
)

type Storage interface {
	GetTasks() ([]model.Task, error)
	GetTasksByDate(date string) ([]model.Task, error)
	GetTasksByString(searchFor string) ([]model.Task, error)
	GetTaskById(id uint64) (model.Task, error)
	InsertTask(task model.Task) (int64, error)
	UpdateTask(task model.Task) error
	DeleteTask(id uint64) error
}

type TodoList struct {
	storage Storage
}

func New(storage Storage) *TodoList {
	return &TodoList{
		storage: storage,
	}
}

func NextDate(date, now time.Time, repeat string) (time.Time, error) {
	if repeat == "" {
		return date, errors.New("no repeat")
	}

	nextDate := date

	repeatParts := strings.Split(repeat, " ")
	switch repeatParts[0] {
	case yearParam:
		for {
			nextDate = nextDate.AddDate(1, 0, 0)

			if nextDate.After(now) {
				break
			}
		}
	case daysParam:
		if len(repeatParts) != 2 {
			return nextDate, errors.New("wrong repeat")
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil {
			return nextDate, err
		}

		if days > daysLimit {
			return nextDate, errors.New("too many days to add")
		}

		for {
			nextDate = nextDate.AddDate(0, 0, days)

			if nextDate.After(now) {
				break
			}
		}
	case weekParam:
		if len(repeatParts) != 2 {
			return nextDate, errors.New("wrong repeat")
		}

		weekdaysParts := strings.Split(repeatParts[1], ",")

		weekdays := make([]int, len(weekdaysParts))
		for i, day := range weekdaysParts {
			weekday, err := strconv.Atoi(day)
			if err != nil {
				return nextDate, err
			}

			if weekday < 1 || weekday > 7 {
				return nextDate, errors.New("wrong weekday")
			}

			weekdays[i] = weekday
		}

		tempNextDate := nextDate

		for {
			tempNextDate = tempNextDate.AddDate(0, 0, 1)

			for _, weekday := range weekdays {
				isBothSunday := weekday == 7 && tempNextDate.Weekday() == 0

				if int(tempNextDate.Weekday()) == (weekday) || isBothSunday {
					nextDate = tempNextDate

					break
				}
			}

			if nextDate.After(now) {
				break
			}
		}
	case monthParam:
		if len(repeatParts) < 2 {
			return nextDate, errors.New("wrong repeat")
		}

		daysParts := strings.Split(repeatParts[1], ",")

		days := make([]int, len(daysParts))
		for i, day := range daysParts {
			dayNum, err := strconv.Atoi(day)
			if err != nil {
				return nextDate, err
			}

			if dayNum < -2 || dayNum > 31 || dayNum == 0 {
				return nextDate, errors.New("wrong day of month")
			}

			days[i] = dayNum
		}

		slices.Sort(days)

		months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		if len(repeatParts) == 3 {
			monthsParts := strings.Split(repeatParts[2], ",")

			months = make([]int, len(monthsParts))
			for i, month := range monthsParts {
				monthNum, err := strconv.Atoi(month)
				if err != nil {
					return nextDate, err
				}

				if monthNum < 1 || monthNum > 12 {
					return nextDate, errors.New("wrong month")
				}

				months[i] = monthNum
			}
		}

		slices.Sort(months)

		year, _, _ := nextDate.Date()
		breaker := false

		for !breaker {
			for _, monthNum := range months {
				for _, dayNum := range days {
					daySetter := dayNum
					if daySetter < 0 {
						daySetter++
					}

					tempDate := time.Date(year, time.Month(monthNum), daySetter, 0, 0, 0, 0, nextDate.Location())

					if tempDate.After(now) && tempDate.After(date) && (tempDate.Month() == time.Month(monthNum) || tempDate.Month() == time.Month(monthNum)-1) {
						nextDate = tempDate
						breaker = true

						break
					}
				}

				if breaker {
					break
				}
			}

			year++
		}
	default:
		return nextDate, errors.New("wrong repeat")
	}

	return nextDate, nil
}

func (t *TodoList) AddTask(task model.Task) (int64, error) {
	return t.storage.InsertTask(task)
}

func (t *TodoList) GetTasks() ([]model.Task, error) {
	return t.storage.GetTasks()
}

func (t *TodoList) GetTasksByDate(date string) ([]model.Task, error) {
	return t.storage.GetTasksByDate(date)
}

func (t *TodoList) GetTasksByString(searchFor string) ([]model.Task, error) {
	return t.storage.GetTasksByString(searchFor)
}

func (t *TodoList) GetTaskById(id uint64) (model.Task, error) {
	return t.storage.GetTaskById(id)
}

func (t *TodoList) UpdateTask(task model.Task) error {
	return t.storage.UpdateTask(task)
}

func (t *TodoList) TaskDone(id int) error {
	task, err := t.storage.GetTaskById(uint64(id))
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return t.storage.DeleteTask(task.ID)
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format in db")
	}

	nextDate, err := NextDate(date, time.Now(), task.Repeat)

	task.Date = nextDate.Format("20060102")

	err = t.storage.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoList) DeleteTask(id int) error {
	return t.storage.DeleteTask(uint64(id))
}
