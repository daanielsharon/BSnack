package commonrouter

import (
	"net/http"
	"time"

	shared_middleware "bsnack/app/pkg/middleware"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func GetParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// DebugMiddleware logs all incoming requests
// func DebugMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("\n=== New Request ===")
// 		log.Printf("Method: %s", r.Method)
// 		log.Printf("URL: %s", r.URL.String())
// 		log.Printf("Path: %s", r.URL.Path)
// 		log.Printf("Headers: %v", r.Header)
// 		log.Printf("Query Params: %v", r.URL.Query())
// 		next.ServeHTTP(w, r)
// 	})
// }

func router() chi.Router {
	r := chi.NewRouter()
	// r.Use(DebugMiddleware)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(shared_middleware.JSONContentType)
	r.Use(shared_middleware.JSONOnly)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(httprate.LimitByIP(2, 1*time.Second))

	return r
}

func New() chi.Router {
	r := router()
	return r
}
