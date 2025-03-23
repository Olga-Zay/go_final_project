package database

import (
	"database/sql"
	"fmt"
	"go_final_project/service/model"
	"strings"
	"time"
)

type DBStorage struct {
	Client *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{Client: db}
}

func (db *DBStorage) CreateTableScheduler() error {
	createTableScheduler := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(256) NOT NULL DEFAULT "",
			comment VARCHAR(256) NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
		);`

	_, err := db.Client.Exec(createTableScheduler)
	if err != nil {
		return fmt.Errorf("Ошибка создания таблицы scheduler в базе данных: %s", err)
	}

	createIndexColumnDate := "CREATE INDEX scheduler_date ON scheduler (date);"

	_, err = db.Client.Exec(createIndexColumnDate)
	if err != nil {
		return fmt.Errorf("Ошибка создания таблицы индекса для колонки date: %s", err)
	}

	return nil
}

func (db *DBStorage) AddTask(taskToAdd Task) (Task, error) {
	addTaskSQL := `INSERT INTO scheduler (
		date, title, comment, repeat
		) VALUES (
		?, ?, ?, ?
	);`

	addingRes, errRes := db.Client.Exec(addTaskSQL, taskToAdd.Date, taskToAdd.Title, taskToAdd.Comment, taskToAdd.Repeat)
	if errRes != nil {
		return Task{}, fmt.Errorf("ошибка сохранения задания в таблице scheduler: %s", errRes)
	}
	taskId, err := addingRes.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf("не удалось получить Id созданного задания: %s", err)
	}
	taskToAdd.Id = int(taskId)

	return taskToAdd, nil
}

func (db *DBStorage) PutTask(taskToSave Task) error {
	putTaskSQL := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?;`

	updateRes, errRes := db.Client.Exec(putTaskSQL, taskToSave.Date, taskToSave.Title, taskToSave.Comment, taskToSave.Repeat, taskToSave.Id)
	if errRes != nil {
		return fmt.Errorf("ошибка сохранения задания в таблице scheduler: %s", errRes.Error())
	}
	rowsUpdated, err := updateRes.RowsAffected()
	if rowsUpdated == 0 {
		if err != nil {
			return fmt.Errorf("не удалось обновить запись с ID %d: %s", taskToSave.Id, err.Error())
		}
		return fmt.Errorf("не удалось обновить запись с ID %d", taskToSave.Id)
	}

	return nil
}

func (db *DBStorage) GetTasks(searchTitle string, searchDate time.Time) ([]Task, error) {
	var partsForWhere []string
	var binds []any
	var dateStr string

	if searchTitle != "" {
		partsForWhere = append(partsForWhere, "title LIKE ?")
		binds = append(binds, "%"+searchTitle+"%")
	}

	if !searchDate.IsZero() {
		dateStr = searchDate.Format(model.CommonDateFormat)
		partsForWhere = append(partsForWhere, "date = ?")
		binds = append(binds, dateStr)
	}

	whereSQL := strings.Join(partsForWhere, " AND ")
	if whereSQL != "" {
		whereSQL = fmt.Sprintf("WHERE %s", whereSQL)
	}

	var tasks []Task
	getTasksSQL := fmt.Sprintf("SELECT * FROM scheduler %s ORDER BY date ASC LIMIT ?;", whereSQL)
	binds = append(binds, model.LimitTasks)
	rows, err := db.Client.Query(getTasksSQL, binds...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *DBStorage) GetTask(id string) (Task, error) {
	var task Task

	getTasksSQL := "SELECT * FROM scheduler WHERE id = ?;"

	err := db.Client.QueryRow(getTasksSQL, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("задача с ID %s не найдена", id)
		}
		return Task{}, err
	}
	return task, nil
}

func (db *DBStorage) DeleteTask(id string) error {
	deleteTaskSQL := "DELETE FROM scheduler WHERE id = ?;"

	deleteRes, errRes := db.Client.Exec(deleteTaskSQL, id)
	if errRes != nil {
		return fmt.Errorf("ошибка удаления задания в таблице scheduler: %s", errRes.Error())
	}
	rowsUpdated, err := deleteRes.RowsAffected()
	if rowsUpdated == 0 {
		if err != nil {
			return fmt.Errorf("не удалось удалить запись с ID %s из-за ошибки %s", id, err.Error())
		}
		return fmt.Errorf("не удалось обновить запись с ID %s", id)
	}

	return nil
}
