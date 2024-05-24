package main

import (
	"log/slog"
	"os"

	"log"

	"github.com/Naumovets/go-search/internal/db/postgres"
	"github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/manager"
	"github.com/Naumovets/go-search/internal/repositories/tasks"
	"github.com/Naumovets/go-search/internal/site"
	_ "github.com/lib/pq"
)

func main() {

	//TODO:
	//rework no complete tasks with status (1)

	cfg_queue, err := postgres.NewConfig(".tasks.env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	loggerCfg, err := logger.NewConfig(".env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	logger.Log = logger.SetupLogger(loggerCfg.Env)
	logger.Info("Starting search-robot", slog.String("env", loggerCfg.Env))

	db, err := postgres.NewConn(*cfg_queue)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	rep := tasks.NewRepository(db)

	exists, err := rep.ExistActualTask()

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	if !exists {
		first_site, err := site.NewSite("https://ru.wikipedia.org/wiki/Заглавная_страница")

		if err != nil {
			os.Exit(1)
		}

		err = rep.AddTask([]site.Site{*first_site})

		if err != nil {
			log.Fatalf("err: %s", err)
			os.Exit(2)
		}
	}

	manager := manager.NewManager(rep)

	manager.Start(100)

	// count, err := rep.GetCountCompleteTasks()

	// if err != nil {
	// 	logger.Debug("err", sl.Err(err))
	// } else {
	// 	logger.Info(fmt.Sprintf("%d", count))
	// }
}
