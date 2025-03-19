package service

import (
	"errors"
	"fmt"
	"go_final_project/database"
	"go_final_project/service/model"
	"strconv"
	"strings"
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
			return model.AddTaskResponse{}, fmt.Errorf("ошибка парсинга даты задачи в AddTask: %s", err)
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
				return model.AddTaskResponse{}, fmt.Errorf("ошибка вычисления следующей даты для просроченной задачи в AddTask: %s", err)
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
		return model.AddTaskResponse{}, fmt.Errorf("ошибка добавления задачи в базу данных: %s", addingErr)
	}

	return model.AddTaskResponse{
		ID: addedTask.Id,
	}, nil
}

func DateParse(dateStr string) (time.Time, error) {
	date, err := time.Parse(model.CommonDateFormat, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, err
}

func PrepareRepeatRuleFromRawString(repeatRuleRaw string) (model.RepeatRule, error) {
	if repeatRuleRaw == "" {
		return model.RepeatRule{}, nil
	}

	repeatSlice := strings.Split(repeatRuleRaw, " ")
	if len(repeatSlice) < 1 {
		return model.RepeatRule{}, errors.New("формат правила повторения не соблюден")
	}

	var repeatValue *int
	if len(repeatSlice) >= 2 {
		rValInt, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return model.RepeatRule{}, fmt.Errorf("не удалось распарсить второй параметр правила повторения: %s", err)
		}
		repeatValue = &rValInt
	}

	return model.RepeatRule{
		Name:  repeatSlice[0],
		Value: repeatValue,
	}, nil
}

func (s *Service) GetClosestTasks() ([]model.Task, error) {
	tasks, err := s.storage.GetTasks()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список задач из базы данных: %s", err)
	}

	return tasks, nil
}
