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

func get2LastMonthDays(date time.Time) (int, int) {
	nextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())

	return nextMonth.AddDate(0, 0, -1).Day(), nextMonth.AddDate(0, 0, -2).Day()
}

func PrepareRepeatRuleFromRawString(repeatRuleRaw string) (model.RepeatRule, error) {
	repeatRule := model.RepeatRule{}
	if repeatRuleRaw == "" {
		return repeatRule, nil
	}

	repeatSlice := strings.Split(repeatRuleRaw, " ")
	rLen := len(repeatSlice)
	if rLen == 0 {
		return repeatRule, errors.New("формат правила повторения не соблюден")
	}

	if rLen > 0 {
		repeatRule.Name = repeatSlice[0]
	}

	if rLen == 1 {
		return repeatRule, nil
	}

	repeatValues := make([][]int, 0, len(repeatSlice)-1)
	if rLen > 1 {
		rValsInts, err := parseRepeatValuesFromString(repeatSlice[1])
		if err != nil {
			return repeatRule, fmt.Errorf("не удалось распарсить 1ю группу значений для правила повторения: %s", err.Error())
		}
		repeatValues = append(repeatValues, rValsInts)
	}

	if rLen > 2 {
		rValsInts, err := parseRepeatValuesFromString(repeatSlice[2])
		if err != nil {
			return repeatRule, fmt.Errorf("не удалось распарсить 2ю группу значений для правила повторения: %s", err.Error())
		}
		repeatValues = append(repeatValues, rValsInts)
	}

	repeatRule.Values = repeatValues

	return repeatRule, nil
}

func parseRepeatValuesFromString(rValsString string) ([]int, error) {
	rValsSlice := strings.Split(rValsString, ",")
	rVals := make([]int, 0, len(rValsSlice))

	for _, rValStr := range rValsSlice {
		rValInt, convErr := strconv.Atoi(rValStr)
		if convErr != nil {
			return rVals, fmt.Errorf("не удалось распарсить значение для правила повторения: %s", convErr.Error())
		}
		rVals = append(rVals, rValInt)
	}

	return rVals, nil
}
