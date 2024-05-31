package main

import (
	"log/slog"
	"os"

	"log"

	"github.com/Naumovets/go-search/internal/db/postgres"
	"github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/manager"
	"github.com/Naumovets/go-search/internal/repositories/db"
	"github.com/Naumovets/go-search/internal/repositories/tasks"
	"github.com/Naumovets/go-search/internal/site"
	_ "github.com/lib/pq"
)

func main() {

	// setting logger
	loggerCfg, err := logger.NewConfig(".env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	logger.Log = logger.SetupLogger(loggerCfg.Env)
	logger.Info("Starting search-robot", slog.String("env", loggerCfg.Env))

	// // setting tasks
	cfgQueue, err := postgres.NewConfig(".tasks.env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	dbTask, err := postgres.NewConn(cfgQueue)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	taskRep := tasks.NewRepository(dbTask)

	exists, err := taskRep.ExistActualTask()

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	if !exists {
		firstSite, err := site.NewSite("https://ru.wikipedia.org/wiki/Заглавная_страница")

		if err != nil {
			os.Exit(1)
		}

		err = taskRep.AddTask([]site.Site{*firstSite})

		if err != nil {
			log.Fatalf("err: %s", err)
			os.Exit(2)
		}
	}

	// setting db
	cfgDB, err := postgres.NewConfig(".db.env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	dbConn, err := postgres.NewConn(cfgDB)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	dbRep := db.NewRepository(dbConn)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(2)
	}

	manager := manager.NewManager(taskRep, dbRep)

	manager.Start(100)

}
