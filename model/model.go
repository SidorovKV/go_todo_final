package model

type Task struct {
	ID      uint64 `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

type TaskDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ResponseTasks struct {
	Tasks []TaskDTO `json:"tasks"`
}

type Password struct {
	Password string `json:"password"`
}
