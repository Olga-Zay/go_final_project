package application

import (
	"github.com/go-chi/chi/v5"
	"go_final_project/application/auth"
	"go_final_project/application/handler"
	"go_final_project/config"
	"log"
	"net/http"
)

type Application struct {
	handler *handler.SchedulerHandler
	config  *config.Config
}

func NewApplication(handler *handler.SchedulerHandler, config *config.Config) Application {
	return Application{handler: handler, config: config}
}

func (a *Application) Start() {
	auth := auth.NewAuth(a.config)

	r := chi.NewRouter()

	r.Use(auth.Middleware)

	r.Handle("/*", http.FileServer(http.Dir("web")))

	r.Get("/api/nextdate", a.handler.NextDate)
	r.Get("/api/task", a.handler.GetTask)
	r.Post("/api/task", a.handler.AddTask)
	r.Put("/api/task", a.handler.PutTask)
	r.Get("/api/tasks", a.handler.GetClosestTasks)
	r.Post("/api/task/done", a.handler.DoTask)
	r.Delete("/api/task", a.handler.DeleteTask)

	r.Post("/api/signin", auth.SingIn)

	svr := &http.Server{
		Addr:    ":" + a.config.Port,
		Handler: r,
	}

	if err := svr.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}
