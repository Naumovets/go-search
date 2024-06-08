package main

import (
	"log"
	"os"

	"github.com/Naumovets/go-search/internal/db/postgres"
	"github.com/Naumovets/go-search/internal/http/application"
	"github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/repositories/db"
	_ "github.com/lib/pq"
)

func main() {

	// setting db
	cfg, err := postgres.NewConfig(".db.env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	conn, err := postgres.NewConn(cfg)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(1)
	}

	rep := db.NewRepository(conn)

	if err != nil {
		log.Fatalf("err: %s\n", err)
		os.Exit(2)
	}

	app := application.New(rep)

	// setting logger
	loggerCfg, err := logger.NewConfig(".env")

	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	logger.Log = logger.SetupLogger(loggerCfg.Env)

	logger.Info("Starting app!")

	err = app.Start()
	if err != nil {
		log.Fatalf("failed to start app: %s", err)
	}
}
