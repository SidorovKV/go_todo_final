package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go_todo_final/config"
	"go_todo_final/model"
	"go_todo_final/services/todolist"
	"go_todo_final/transform"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	webDir = "./web"

	nowParam    = "now"
	dateParam   = "date"
	idParam     = "id"
	repeatParam = "repeat"
	searchParam = "search"
)

type TodoList interface {
	AddTask(task model.Task) (int64, error)
	GetTasks() ([]model.Task, error)
	GetTasksByDate(date string) ([]model.Task, error)
	GetTasksByString(searchFor string) ([]model.Task, error)
	GetTaskById(id uint64) (model.Task, error)
	UpdateTask(task model.Task) error
	TaskDone(id int) error
	DeleteTask(id int) error
}

type Server struct {
	httpServer *http.Server
	todoList   TodoList
}

func New(cfg *config.Config, todoList TodoList) *Server {
	s := &Server{
		httpServer: &http.Server{
			Addr:              cfg.Port,
			ReadHeaderTimeout: 1 * time.Second,
		},
		todoList: todoList,
	}

	s.setupHandlers()

	return s
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) setupHandlers() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("GET /api/nextdate", s.nextDateHandler)
	http.HandleFunc("POST /api/task/done", s.taskDoneHandler)
	http.HandleFunc("DELETE /api/task", s.deleteTaskHandler)
	http.HandleFunc("POST /api/task", s.addTaskHandler)
	http.HandleFunc("GET /api/task", s.getTaskByIdHandler)
	http.HandleFunc("PUT /api/task", s.updateTaskHandler)
	http.HandleFunc("GET /api/tasks", s.getTasksHandler)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) nextDateHandler(w http.ResponseWriter, r *http.Request) {
	date, err := time.Parse("20060102", r.URL.Query().Get(dateParam))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	now, err := time.Parse("20060102", r.URL.Query().Get(nowParam))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	repeat := r.URL.Query().Get(repeatParam)
	nextDate, err := todolist.NextDate(date, now, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(nextDate.Format("20060102")))
}

func (s *Server) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var incomingTask model.TaskDTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&incomingTask); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("json Decoder:", err)

		return
	}

	task, err := transform.DtoToTask(incomingTask)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("DtoToTask:", err)

		return
	}

	id, err := s.todoList.AddTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("AddTask:", err)

		return
	}

	w.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, id)))

}

func (s *Server) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var (
		tasks []model.Task
		err   error
	)

	if r.URL.Query().Has(searchParam) {
		tasks, err = s.search(r.URL.Query().Get(searchParam))
	} else {
		tasks, err = s.todoList.GetTasks()
	}

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("GetTasks:", err)

		return
	}

	dtos := transform.TasksToDto(tasks)

	response := model.ResponseTasks{Tasks: dtos}

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}

	w.Write(responseBody)
}

func (s *Server) getTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(r.URL.Query().Get(idParam))
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	task, err := s.todoList.GetTaskById(uint64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	dto := transform.TaskToDto(task)

	responseBody, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}

	w.Write(responseBody)
}

func (s *Server) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var incomingTask model.TaskDTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&incomingTask); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("json Decoder:", err)

		return
	}

	task, err := transform.DtoToTask(incomingTask)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("DtoToTask:", err)

		return
	}

	if task.ID == 0 {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)

		return
	}

	err = s.todoList.UpdateTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("UpdateTask:", err)

		return
	}

	w.Write([]byte(`{}`))
}

func (s *Server) taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(idParam))
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	err = s.todoList.TaskDone(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(`{}`))
}

func (s *Server) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(idParam))
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	err = s.todoList.DeleteTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(`{}`))
}

func (s *Server) search(searchQuery string) ([]model.Task, error) {
	date, err := time.Parse("02.01.2006", searchQuery)
	if err != nil {
		return s.todoList.GetTasksByString(searchQuery)
	}

	return s.todoList.GetTasksByDate(date.Format("20060102"))
}
