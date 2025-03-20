package application

import (
	"github.com/go-chi/chi/v5"
	"go_final_project/application/handler"
	"log"
	"net/http"
)

type Application struct {
	handler *handler.SchedulerHandler
	srvPort string
}

func NewApplication(handler *handler.SchedulerHandler, srvPort string) Application {
	return Application{handler: handler, srvPort: srvPort}
}

func (a *Application) Start() {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", a.handler.NextDate)
	r.Get("/api/task", a.handler.GetTask)
	r.Post("/api/task", a.handler.AddTask)
	r.Put("/api/task", a.handler.PutTask)
	r.Get("/api/tasks", a.handler.GetClosestTasks)
	r.Post("/api/task/done", a.handler.DoTask)
	r.Delete("/api/task", a.handler.DeleteTask)

	svr := &http.Server{
		Addr:    ":" + a.srvPort,
		Handler: r,
	}

	if err := svr.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}
