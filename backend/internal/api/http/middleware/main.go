package middleware

type Middleware interface {
}

type middleware struct{}

func New() Middleware {
	return &middleware{}
}
