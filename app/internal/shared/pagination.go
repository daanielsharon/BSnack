package shared

import (
	ctxkey "bsnack/app/pkg/ctxKey"
	"bsnack/app/pkg/middleware"
	"context"
)

func GetPagination(ctx context.Context) *middleware.Pagination {
	pg, ok := ctx.Value(ctxkey.PaginationKey()).(*middleware.Pagination)
	if !ok {
		return &middleware.Pagination{
			Page:    1,
			PerPage: 10,
			Offset:  0,
		}
	}
	return pg
}
