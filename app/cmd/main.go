package main

import (
	"bsnack/app/internal/router"
	"bsnack/app/pkg/common"
	"bsnack/app/pkg/config"
	"bsnack/app/pkg/db"
	"bsnack/app/pkg/envLoader"
	"bsnack/app/pkg/validation"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigration(dsn string) {
	wd, _ := os.Getwd()
	m, err := migrate.New(
		"file://"+wd+"/app/migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("Migration init error: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migration needed.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Println("Migration completed successfully.")
	}
}

func main() {
	envLoader.LoadEnv(".env")
	cfg := config.LoadConfig("config.yaml")

	if os.Getenv("MIGRATE") == "true" {
		runMigration(common.GenerateDBUrl(cfg))
		os.Exit(0)
	}

	dsn := common.GenerateDSN(cfg)
	port := cfg.Server.Port

	validation.Init()
	DB, err := db.Connect(db.Config{
		DSN: dsn,
	})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	redisClient := db.NewRedisClient(db.RedisConfig{
		Address: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	routes := router.NewRouter(DB, redisClient)
	log.Println("Server started on", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
}
