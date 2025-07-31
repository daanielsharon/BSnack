package middleware

import (
	ctxkey "bsnack/app/pkg/ctxKey"
	"context"
	"net/http"
	"strconv"
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

func GetPagination(ctx context.Context) *Pagination {
	pg, ok := ctx.Value(ctxkey.PaginationKey()).(*Pagination)
	if !ok {
		return &Pagination{
			Page:    1,
			PerPage: 10,
			Offset:  0,
		}
	}
	return pg
}

func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		page, _ := strconv.Atoi(q.Get("page"))
		if page < 1 {
			page = 1
		}

		perPage, _ := strconv.Atoi(q.Get("limit"))
		if perPage < 1 {
			perPage = 10
		}

		offset := (page - 1) * perPage

		ctx := context.WithValue(r.Context(), ctxkey.PaginationKey(), &Pagination{
			Page:    page,
			PerPage: perPage,
			Offset:  offset,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
