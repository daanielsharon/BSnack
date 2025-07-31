package ctxkey

type contextKey string

const paginationKey contextKey = "pagination"

func PaginationKey() contextKey {
	return paginationKey
}
