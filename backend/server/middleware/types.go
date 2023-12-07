package middleware

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler
