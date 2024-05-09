package main

import (
	"log/slog"
	"os"

	"log"

	"github.com/Naumovets/go-search/internal/db/postgres"
	"github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/manager"
	"github.com/Naumovets/go-search/internal/repositories/tasks"
	_ "github.com/lib/pq"
)

// TODO:
// design tests: internal/site/site_test.go

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

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

	logger.Log = setupLogger(loggerCfg.Env)
	logger.Info("Starting go-search", slog.String("env", loggerCfg.Env))

	db, err := postgres.NewConn(*cfg_queue)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	rep := tasks.NewRepository(db)

	// site1, err := site.NewSite("https://vk.com/daniilnaumovets")

	// if err != nil {
	// 	os.Exit(1)
	// }

	// err = rep.AddTask([]site.Site{*site1})

	// if err != nil {
	// 	fmt.Printf("err: %s", err)
	// 	os.Exit(2)
	// }

	manager := manager.NewManager(rep)

	manager.Start(10)

	// res, err := rep.GetLimitTasks(5000000)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(len(res))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
