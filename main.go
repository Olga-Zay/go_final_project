package main

import (
	"database/sql"
	"go_final_project/application"
	"go_final_project/application/handler"
	"go_final_project/config"
	"go_final_project/database"
	"go_final_project/service"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

func main() {
	//Вынести всё в слой ресурсов
	cfg := config.LoadConfig()

	//appPath, err := os.Executable()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	dbFile := cfg.DB
	_, err := os.Stat(cfg.DB)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	defer db.Close()

	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %s", err)
	}

	dbStorage := database.NewDBStorage(db)
	appHandler := handler.NewSchedulerHandler(service.NewService(dbStorage))

	if install {
		err := dbStorage.CreateTableScheduler()
		if err != nil {
			log.Fatalf("Ошибка создания таблицы: %s", err)
		}
	}

	app := application.NewApplication(appHandler, cfg.Port)
	app.Start()
}
