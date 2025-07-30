package router

import (
	v1 "bsnack/app/internal/api/v1"
	commonrouter "bsnack/app/pkg/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) chi.Router {
	r := commonrouter.New()

	r.Route("/api", func(r chi.Router) {
		r.Mount("/v1", v1.Routes(v1.Router{}))
	})

	return r
}
