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
)

func main() {
	var dsn string
	var port string

	envLoader.LoadEnv(".env")
	cfg := config.LoadConfig("config.yaml")
	dsn = common.GenerateDSN(*cfg)
	port = cfg.Server.Port

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
