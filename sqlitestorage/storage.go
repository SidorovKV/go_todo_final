package sqlitestorage

import (
	"fmt"
	"go_todo_final/config"
	"go_todo_final/model"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbLimit = 50

type Storage struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (*Storage, error) {
	db, err := sqlx.Connect("sqlite3", cfg.DBFile)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	s := &Storage{db: db}

	err = s.tryCreateTable()
	if err != nil {
		return nil, fmt.Errorf("tryCreateTable: %w", err)
	}

	err = os.Chmod(cfg.DBFile, 0777)
	if err != nil {
		return nil, fmt.Errorf("chmod: %w", err)
	}

	return s, nil
}

// Бесполезная проверка, как по мне. Даже если файл существует, нет гарантии, что там существует нужная таблица.
func checkDbFileExists(cfg *config.Config) bool {
	_, err := os.Stat(cfg.DBFile)
	if err != nil {
		return false
	}

	return true
}

func (s *Storage) tryCreateTable() error {

	_, err := s.db.Exec(createTablequery)
	if err != nil {
		return fmt.Errorf("table create: %w", err)
	}

	_, err = s.db.Exec(createIndexQuery)
	if err != nil {
		return fmt.Errorf("index create: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) InsertTask(task model.Task) (int64, error) {
	res, err := s.db.Exec(insertTaskQuery, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("insertTask failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting LastInsertId failed: %w", err)
	}
	return id, nil
}

func (s *Storage) GetTasks() ([]model.Task, error) {
	var tasks []model.Task

	err := s.db.Select(&tasks, getTasksQuery, time.Now().Format("20060102"), dbLimit)
	if err != nil {
		return nil, fmt.Errorf("selectTasks failed: %w", err)
	}

	return tasks, nil
}

func (s *Storage) GetTasksByDate(date string) ([]model.Task, error) {
	var tasks []model.Task

	err := s.db.Select(&tasks, getTasksByDateQuery, date, dbLimit)
	if err != nil {
		return nil, fmt.Errorf("selectTasks failed: %w", err)
	}

	return tasks, nil
}

func (s *Storage) GetTasksByString(searchFor string) ([]model.Task, error) {
	var tasks []model.Task

	searchFor = "%" + searchFor + "%"
	err := s.db.Select(&tasks, getTasksByStringQuery, searchFor, searchFor, dbLimit)
	if err != nil {
		return nil, fmt.Errorf("selectTasks failed: %w", err)
	}

	return tasks, nil
}

func (s *Storage) GetTaskById(id uint64) (model.Task, error) {
	var task model.Task

	err := s.db.Get(&task, getTaskByIdQuery, id)
	if err != nil {
		return model.Task{}, fmt.Errorf("selectTask failed: %w", err)
	}

	return task, nil
}

func (s *Storage) UpdateTask(task model.Task) error {
	res, err := s.db.Exec(updateTaskQuery, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("updateTask failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting RowsAffected failed: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("update failed")
	}

	return nil
}

func (s *Storage) DeleteTask(id uint64) error {
	res, err := s.db.Exec(deleteTaskQuery, id)
	if err != nil {
		return fmt.Errorf("deleteTask failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting RowsAffected failed: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("delete failed")
	}

	return nil
}
