package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/application/validator"
	"go_final_project/service"
	"go_final_project/service/model"
	"net/http"
)

type SchedulerHandler struct {
	service *service.Service
}

func NewSchedulerHandler(service *service.Service) *SchedulerHandler {
	return &SchedulerHandler{service: service}
}

func (h *SchedulerHandler) NextDate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("repeat") == "" {
		//TODO написать sql запрос на удаление задачи из бд
	}

	//Получение и валидация данных
	nextDateRequest, err := h.prepareNextDateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validator.ValidateNextDateRequest(nextDateRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Вызов бизнес логики
	nextDay, err := h.service.CalculateNextDate(nextDateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(nextDay.Format(model.CommonDateFormat)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *SchedulerHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	addTaskRequest, err := h.prepareAddTaskRequest(r)
	if err != nil {
		errResp := &model.AddTaskResponse{
			Error: fmt.Sprintf("валидация запроса не пройдена: %s", err.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	if errValid := validator.ValidateAddTaskRequest(addTaskRequest); errValid != nil {
		errResp := &model.AddTaskResponse{
			Error: fmt.Sprintf("валидация запроса не пройдена: %s", errValid.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	addTaskResponse, respErr := h.service.AddTask(addTaskRequest)
	if respErr != nil {
		addTaskResponse.Error = fmt.Sprintf("ошибка при добавлении задания: %s", respErr.Error())
		h.prepareTaskResponse(w, &addTaskResponse, http.StatusInternalServerError)
		return
	}

	h.prepareTaskResponse(w, &addTaskResponse, http.StatusOK)
}

func (h *SchedulerHandler) GetClosestTasks(w http.ResponseWriter, r *http.Request) {
	tasksResp := model.ClosestTasksResponse{
		Tasks: []model.Task{},
	}
	tasks, err := h.service.GetClosestTasks()
	if err != nil {
		tasksResp.Error = err.Error()
		h.prepareTaskResponse(w, &tasksResp, http.StatusInternalServerError)
		return
	}

	if tasks != nil {
		tasksResp.Tasks = tasks
	}
	h.prepareTaskResponse(w, &tasksResp, http.StatusOK)
}

func (h *SchedulerHandler) prepareNextDateRequest(r *http.Request) (model.NextDateRequest, error) {
	nowStr := r.URL.Query().Get("now")
	dateNow, err := service.DateParse(nowStr)
	if err != nil {
		return model.NextDateRequest{}, err
	}

	date, err := service.DateParse(r.URL.Query().Get("date"))
	if err != nil {
		return model.NextDateRequest{}, err
	}

	repeatStr := r.URL.Query().Get("repeat")

	if repeatStr == "" {
		return model.NextDateRequest{}, errors.New("не указано правило повторения")
	}

	repeatRule, err := service.PrepareRepeatRuleFromRawString(repeatStr)
	if err != nil {
		return model.NextDateRequest{}, fmt.Errorf("ошибка парсига правил повторения при вычислении nextDate: %s", err.Error())
	}

	nextDateRequest := model.NextDateRequest{
		Now:    dateNow,
		Date:   date,
		Repeat: repeatRule,
	}

	return nextDateRequest, nil
}

func (h *SchedulerHandler) prepareAddTaskRequest(r *http.Request) (model.AddTaskRequest, error) {
	var addTaskRequest model.AddTaskRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&addTaskRequest); err != nil {
		return model.AddTaskRequest{}, fmt.Errorf("ошибка десериализации JSON: %s", err.Error())
	}

	if addTaskRequest.RepeatRaw != "" {
		repeatRule, err := service.PrepareRepeatRuleFromRawString(addTaskRequest.RepeatRaw)
		if err != nil {
			return model.AddTaskRequest{}, fmt.Errorf("ошибка парсига правил повторения при создании задания: %s", err.Error())
		}
		addTaskRequest.Repeat = repeatRule
	}

	return addTaskRequest, nil
}

func (h *SchedulerHandler) prepareTaskResponse(w http.ResponseWriter, addTaskResponse any, httpStatus int) {
	encoderErr := json.NewEncoder(w).Encode(&addTaskResponse)
	if encoderErr != nil {
		http.Error(w, encoderErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
}
