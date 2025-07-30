package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
}

func Routes(router Router) http.Handler {
	r := chi.NewRouter()

	return r
}
