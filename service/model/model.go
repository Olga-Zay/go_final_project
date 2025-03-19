package model

import "time"

type NextDateRequest struct {
	Now    time.Time
	Date   time.Time
	Repeat RepeatRule
}

type RepeatRule struct {
	Name  string
	Value *int
}

type AddTaskRequest struct {
	Date      string `json:"date"`
	Title     string `json:"title"`
	Comment   string `json:"comment"`
	RepeatRaw string `json:"repeat"`
	Repeat    RepeatRule
}

type AddTaskResponse struct {
	ID    int    `json:"id"`
	Error string `json:"error"`
}

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ClosestTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Error string `json:"error"`
}

type GetTaskRequest struct {
	TaskId string `json:"id"`
}

type GetTaskResponse struct {
	Task
}

type GetTaskResponseWithError struct {
	Error string `json:"error"`
}
