package service

import (
	"errors"
	"fmt"
	"go_final_project/service/model"
	"strconv"
	"strings"
	"time"
)

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
