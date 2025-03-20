package service

import (
	"fmt"
	"go_final_project/database"
	"go_final_project/service/model"
	"strconv"
	"time"
)

type Service struct {
	storage *database.DBStorage
}

func NewService(storage *database.DBStorage) *Service {
	return &Service{
		storage: storage,
	}
}

// CalculateNextDate вычисляет корректную новую дату задания на основе переданного правила повторения
func (s *Service) CalculateNextDate(nextDateRequest model.NextDateRequest) (time.Time, error) {
	var newDate time.Time

	if nextDateRequest.Repeat.Name == "y" {
		newDate = nextDateRequest.Date.AddDate(1, 0, 0)

		for newDate.Before(nextDateRequest.Now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate, nil
	}

	newDate = nextDateRequest.Date.AddDate(0, 0, *nextDateRequest.Repeat.Value)

	for newDate.Before(nextDateRequest.Now) {
		newDate = newDate.AddDate(0, 0, *nextDateRequest.Repeat.Value)
	}

	return newDate, nil
}

// AddTask добавляет задание
func (s *Service) AddTask(addTaskRequest model.AddTaskRequest) (model.AddTaskResponse, error) {
	nowDate := time.Now()
	nowDateStr := nowDate.Format(model.CommonDateFormat)
	taskDate := nowDate
	// Если дата в запросе не указана, то сегодняшнюю берём
	if addTaskRequest.Date != "" {
		reqDate, err := DateParse(addTaskRequest.Date)
		if err != nil {
			return model.AddTaskResponse{}, fmt.Errorf("ошибка парсинга даты задачи в AddTask: %s", err.Error())
		}

		//если правило повторения не указано, продолжаем с сегодняшним числом
		if addTaskRequest.Date < nowDateStr && addTaskRequest.RepeatRaw != "" {
			// при указанном правиле повторения вычислем новую дату выполнения,
			nextDate, nextDateErr := s.CalculateNextDate(model.NextDateRequest{
				Now:    nowDate,
				Date:   reqDate,
				Repeat: addTaskRequest.Repeat,
			})
			if nextDateErr != nil {
				return model.AddTaskResponse{}, fmt.Errorf("ошибка вычисления следующей даты для просроченной задачи в AddTask: %s", err.Error())
			}
			taskDate = nextDate
		}
	}

	addedTask, addingErr := s.storage.AddTask(database.Task{
		Date:    taskDate.Format(model.CommonDateFormat),
		Title:   addTaskRequest.Title,
		Comment: addTaskRequest.Comment,
		Repeat:  addTaskRequest.RepeatRaw,
	})
	if addingErr != nil {
		return model.AddTaskResponse{}, fmt.Errorf("ошибка добавления задачи в базу данных: %s", addingErr.Error())
	}

	return model.AddTaskResponse{
		ID: addedTask.Id,
	}, nil
}

// GetTask получить запрошенное задание
func (s *Service) GetTask(request model.GetTaskRequest) (model.Task, error) {
	task, err := s.storage.GetTask(request.TaskId)
	if err != nil {
		return model.Task{}, fmt.Errorf("ошибка получения задачи из базы данных: %s", err.Error())
	}

	return model.Task{
		Id:      strconv.Itoa(task.Id),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}, nil
}

// DoTask получить запрошенное задание
func (s *Service) DoTask(request model.DoTaskRequest, onlyDelete bool) (bool, error) {
	taskToBeDone, err := s.storage.GetTask(request.TaskId)
	if err != nil {
		return false, fmt.Errorf("не удалось получить задачу для выполнения: %s", err.Error())
	}
	if onlyDelete || taskToBeDone.Repeat == "" {
		deleteErr := s.storage.DeleteTask(request.TaskId)
		if err != nil {
			return false, fmt.Errorf("не удалось удалить задачу из базы данных: %s", deleteErr.Error())
		}

		return true, nil
	}

	prevTaskDate, err := DateParse(taskToBeDone.Date)
	if err != nil {
		return false, fmt.Errorf("не удалось вычислить дату следующего выполнения: %s", err.Error())
	}
	repeatRule, err := PrepareRepeatRuleFromRawString(taskToBeDone.Repeat)
	if err != nil {
		return false, fmt.Errorf("не удалось вычислить дату следующего выполнения: %s", err.Error())
	}
	nextDate, nextDateErr := s.CalculateNextDate(model.NextDateRequest{
		Now:    time.Now(),
		Date:   prevTaskDate,
		Repeat: repeatRule,
	})
	if nextDateErr != nil {
		return false, fmt.Errorf("не удалось вычислить дату следующего выполнения: %s", nextDateErr.Error())
	}
	taskToBeDone.Date = nextDate.Format(model.CommonDateFormat)

	editErr := s.storage.PutTask(taskToBeDone)
	if editErr != nil {
		return false, fmt.Errorf("ошибка редактирования выполняемой задачи в базе данных: %s", editErr.Error())
	}

	return true, nil
}

// PutTask отредактировать информацию задания
func (s *Service) PutTask(request model.PutTaskRequest) (bool, error) {
	reqDate, err := DateParse(request.Date)
	if err != nil {
		return false, fmt.Errorf("ошибка парсинга даты задачи в PutTask: %s", err.Error())
	}
	reqDate = reqDate.AddDate(0, 0, 1)
	now := time.Now()

	// Если дата в запросе не указана или меньше сегодняшней, то ошибка
	if request.Date == "" || reqDate.Before(now) {
		return false, fmt.Errorf("дата задания указана неверно для PutTask: %s", request.Date)
	}

	taskId, convErr := strconv.Atoi(request.Id)
	if convErr != nil {
		return false, fmt.Errorf("передан не числовой ID задания: %s", convErr.Error())
	}

	editErr := s.storage.PutTask(database.Task{
		Id:      taskId,
		Date:    request.Date,
		Title:   request.Title,
		Comment: request.Comment,
		Repeat:  request.Repeat,
	})
	if editErr != nil {
		return false, fmt.Errorf("ошибка редактирования задачи в базе данных: %s", editErr.Error())
	}

	return true, nil
}

// GetClosestTasks получить ближайшие задачи
func (s *Service) GetClosestTasks() ([]model.Task, error) {
	dbTasks, err := s.storage.GetTasks()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список задач из базы данных: %s", err.Error())
	}

	tasks := make([]model.Task, 0, len(dbTasks))
	for _, task := range dbTasks {
		tasks = append(tasks, model.Task{
			Id:      strconv.Itoa(task.Id),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	return tasks, nil
}
