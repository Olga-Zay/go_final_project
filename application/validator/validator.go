package validator

import (
	"errors"
	"fmt"
	"go_final_project/service"
	"go_final_project/service/model"
)

//Было бы хорошо добавить здесь объект валидатора и передавать его в качестве зависимости тем, кому нужно

func ValidateRepeat(repeat model.RepeatRule) error {
	if repeat.Name == "" {
		return errors.New("формат правила повторения не соблюден")
	}

	if repeat.Name != "d" && repeat.Name != "y" {
		return errors.New("формат правила повторения не соблюден")
	}

	rVal := repeat.Value

	if rVal != nil {
		if *rVal > 400 || *rVal < 1 {
			return errors.New("формат правила повторения не соблюден")
		}
	} else if repeat.Name != "y" {
		return errors.New("формат правила повторения не соблюден")
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
