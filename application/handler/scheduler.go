package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/service"
	"go_final_project/service/model"
	"go_final_project/service/validator"
	"net/http"
	"time"
)

type SchedulerHandler struct {
	service *service.Service
}

func NewSchedulerHandler(service *service.Service) *SchedulerHandler {
	return &SchedulerHandler{service: service}
}

func (h *SchedulerHandler) NextDate(w http.ResponseWriter, r *http.Request) {
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
			Error: fmt.Sprintf("не удалось распарсить данные запроса: %s", err.Error()),
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

func (h *SchedulerHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	request, err := h.prepareGetTaskRequest(r)
	if err != nil {
		errResp := &model.GetTaskResponseWithError{
			Error: fmt.Sprintf("не удалось распарсить данные запроса: %s", err.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	if errValid := validator.ValidateGetTaskRequest(request); errValid != nil {
		errResp := &model.GetTaskResponseWithError{
			Error: fmt.Sprintf("валидация запроса не пройдена: %s", errValid.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	task, serviceErr := h.service.GetTask(request)
	if serviceErr != nil {
		getTaskResponse := model.GetTaskResponseWithError{
			Error: fmt.Sprintf("ошибка при поиске задания: %s", serviceErr.Error()),
		}
		h.prepareTaskResponse(w, &getTaskResponse, http.StatusInternalServerError)
		return
	}

	getTaskResponse := model.GetTaskResponse{Task: task}

	h.prepareTaskResponse(w, &getTaskResponse, http.StatusOK)
}

func (h *SchedulerHandler) PutTask(w http.ResponseWriter, r *http.Request) {
	request, err := h.preparePutTaskRequest(r)
	if err != nil {
		errResp := &model.PutTaskResponseWithError{
			Error: fmt.Sprintf("не удалось распарсить данные запроса: %s", err.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	if errValid := validator.ValidatePutTaskRequest(request); errValid != nil {
		errResp := &model.PutTaskResponseWithError{
			Error: fmt.Sprintf("валидация запроса не пройдена: %s", errValid.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	_, serviceErr := h.service.PutTask(request)
	if serviceErr != nil {
		putTaskResponse := model.PutTaskResponseWithError{
			Error: fmt.Sprintf("ошибка при редактировании задания: %s", serviceErr.Error()),
		}
		h.prepareTaskResponse(w, &putTaskResponse, http.StatusInternalServerError)
		return
	}

	h.prepareTaskResponse(w, &model.PutTaskResponse{}, http.StatusOK)
}

func (h *SchedulerHandler) DoTask(w http.ResponseWriter, r *http.Request) {
	h.doTask(w, r, false)
}

func (h *SchedulerHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.doTask(w, r, true)
}

func (h *SchedulerHandler) GetClosestTasks(w http.ResponseWriter, r *http.Request) {
	closestTasksRequest, err := h.prepareGetClosestTasksRequest(r)
	if err != nil {
		tasksRespErr := model.ClosestTasksResponseWithError{
			Error: fmt.Sprintf("не удалось получить параметры запроса: %s", err.Error()),
		}
		h.prepareTaskResponse(w, &tasksRespErr, http.StatusBadRequest)
		return
	}

	tasks, err := h.service.GetClosestTasks(closestTasksRequest)
	if err != nil {
		tasksRespErr := model.ClosestTasksResponseWithError{
			Error: fmt.Sprintf("не удалось получить ближайшие задачи: %s", err.Error()),
		}
		h.prepareTaskResponse(w, &tasksRespErr, http.StatusInternalServerError)
		return
	}

	h.prepareTaskResponse(w, &model.ClosestTasksResponse{Tasks: tasks}, http.StatusOK)
}

func (h *SchedulerHandler) doTask(w http.ResponseWriter, r *http.Request, onlyDelete bool) {
	request, err := h.prepareDoTaskRequest(r)
	if err != nil {
		errResp := &model.DoTaskResponseWithError{
			Error: fmt.Sprintf("не удалось распарсить данные запроса: %s", err.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	if errValid := validator.ValidateDoTaskRequest(request); errValid != nil {
		errResp := &model.DoTaskResponseWithError{
			Error: fmt.Sprintf("валидация запроса не пройдена: %s", errValid.Error()),
		}
		h.prepareTaskResponse(w, errResp, http.StatusBadRequest)
		return
	}

	_, serviceErr := h.service.DoTask(request, onlyDelete)
	if serviceErr != nil {
		response := model.DoTaskResponseWithError{
			Error: fmt.Sprintf("ошибка при выполнении задания: %s", serviceErr.Error()),
		}
		h.prepareTaskResponse(w, &response, http.StatusInternalServerError)
		return
	}

	h.prepareTaskResponse(w, &model.DoTaskResponse{}, http.StatusOK)
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

func (h *SchedulerHandler) preparePutTaskRequest(r *http.Request) (model.PutTaskRequest, error) {
	var putTaskRequest model.PutTaskRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&putTaskRequest); err != nil {
		return model.PutTaskRequest{}, fmt.Errorf("ошибка десериализации JSON: %s", err.Error())
	}

	if putTaskRequest.Repeat != "" {
		repeatRule, err := service.PrepareRepeatRuleFromRawString(putTaskRequest.Repeat)
		if err != nil {
			return model.PutTaskRequest{}, fmt.Errorf("ошибка парсига правил повторения при редактировании задания: %s", err.Error())
		}

		putTaskRequest.RepeatRule = repeatRule
	}

	return putTaskRequest, nil
}

func (h *SchedulerHandler) prepareGetTaskRequest(r *http.Request) (model.GetTaskRequest, error) {
	return model.GetTaskRequest{
		TaskId: r.URL.Query().Get("id"),
	}, nil
}

func (h *SchedulerHandler) prepareDoTaskRequest(r *http.Request) (model.DoTaskRequest, error) {
	return model.DoTaskRequest{
		TaskId: r.URL.Query().Get("id"),
	}, nil
}

func (h *SchedulerHandler) prepareTaskResponse(w http.ResponseWriter, addTaskResponse any, httpStatus int) {
	encoderErr := json.NewEncoder(w).Encode(&addTaskResponse)
	if encoderErr != nil {
		http.Error(w, encoderErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
}

func (h *SchedulerHandler) prepareGetClosestTasksRequest(r *http.Request) (model.ClosestTasksRequest, error) {
	searchVal := r.URL.Query().Get("search")
	if searchVal == "" {
		return model.ClosestTasksRequest{}, nil
	}

	searchDate, err := time.Parse(model.SearchDateFormat, searchVal)
	if err != nil {
		return model.ClosestTasksRequest{SearchTitle: searchVal}, nil
	}

	return model.ClosestTasksRequest{SearchDate: searchDate}, nil
}
