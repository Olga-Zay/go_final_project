package validator

import (
	"errors"
	"fmt"
	"go_final_project/service"
	"go_final_project/service/model"
)

var ValidRepeatRuleNames = map[string]bool{
	"d": true,
	"y": true,
	"w": true,
	"m": true,
}

func ValidateRepeat(repeat model.RepeatRule) error {
	if _, ok := ValidRepeatRuleNames[repeat.Name]; !ok {
		return errors.New("формат правила повторения не соблюден")
	}

	switch repeat.Name {
	case "d":
		if len(repeat.Values) != 1 || len(repeat.Values[0]) != 1 {
			return errors.New("формат правила повторения для D не соблюден")
		}

		dVal := repeat.Values[0][0]
		if dVal > 400 || dVal < 1 {
			return errors.New("формат правила повторения для D не соблюден")
		}
	case "w":
		if len(repeat.Values) != 1 || len(repeat.Values[0]) == 0 {
			return errors.New("формат правила повторения для W не соблюден")
		}

		for _, wVal := range repeat.Values[0] {
			if wVal < 1 || wVal > 7 {
				return errors.New("формат правила повторения для W не соблюден")
			}
		}
	case "m":
		mValsLen := len(repeat.Values)
		if mValsLen == 0 || mValsLen > 2 ||
			mValsLen == 1 && len(repeat.Values[0]) == 0 ||
			mValsLen == 2 && len(repeat.Values[1]) == 0 {
			return errors.New("формат правила повторения для M не соблюден")
		}

		if len(repeat.Values) == 1 {
			for _, mValDay := range repeat.Values[0] {
				if mValDay == 0 || mValDay < -2 || mValDay > 31 {
					return errors.New("формат правила повторения для M обозначающего номер дня не соблюден")
				}
			}
		}

		if len(repeat.Values) == 2 {
			for _, mValMonth := range repeat.Values[1] {
				if mValMonth < 1 || mValMonth > 12 {
					return errors.New("формат правила повторения для M обозначающего номер месяца не соблюден")
				}
			}
		}
	case "y":
		if len(repeat.Values) > 0 {
			return errors.New("формат правила повторения для Y не соблюден")
		}
	}

	return nil
}

func ValidateNextDateRequest(nextDateRequest model.NextDateRequest) error {
	err := ValidateRepeat(nextDateRequest.Repeat)
	if err != nil {
		return err
	}

	return nil
}

func ValidateAddTaskRequest(addTaskRequest model.AddTaskRequest) error {
	if addTaskRequest.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	if addTaskRequest.Date != "" {
		_, err := service.DateParse(addTaskRequest.Date)
		if err != nil {
			return fmt.Errorf("дата представлена в формате, отличном от %s: %s", model.CommonDateFormat, err.Error())
		}
	}

	if addTaskRequest.RepeatRaw != "" {
		err := ValidateRepeat(addTaskRequest.Repeat)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidatePutTaskRequest(request model.PutTaskRequest) error {
	if request.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	if request.Date != "" {
		_, err := service.DateParse(request.Date)
		if err != nil {
			return fmt.Errorf("дата представлена в формате, отличном от %s: %s", model.CommonDateFormat, err.Error())
		}
	}

	if request.Repeat != "" {
		err := ValidateRepeat(request.RepeatRule)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateGetTaskRequest(request model.GetTaskRequest) error {
	if request.TaskId == "" {
		return errors.New("не указан идентификатор задачи")
	}

	return nil
}

func ValidateDoTaskRequest(request model.DoTaskRequest) error {
	if request.TaskId == "" {
		return errors.New("не указан идентификатор задачи")
	}

	return nil
}
