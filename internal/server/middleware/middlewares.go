package middleware

import "github.com/shanto-323/axis/internal/server"

type Middlewares struct {
	*Global
	*RateLimit
	*ContextEnhancer
	*AuthMiddleware
}

func New(s *server.Server) *Middlewares {
	return &Middlewares{
		Global:          NewGlobal(s),
		RateLimit:       NewRateLimit(s),
		ContextEnhancer: NewContextEnhancer(s),
		AuthMiddleware:  NewAuthMiddleware(s),
	}
}
