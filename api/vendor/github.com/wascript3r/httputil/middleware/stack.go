package middleware

import (
	"github.com/julienschmidt/httprouter"
)

type Middleware func(next httprouter.Handle) httprouter.Handle

type Stack struct {
	middlewares []Middleware
}

func New() *Stack {
	return &Stack{nil}
}

func (s *Stack) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Stack) Wrap(fn httprouter.Handle) httprouter.Handle {
	l := len(s.middlewares)
	if l == 0 {
		return fn
	}

	var result httprouter.Handle
	result = s.middlewares[l-1](fn)

	for i := l - 2; i >= 0; i-- {
		result = s.middlewares[i](result)
	}

	return result
}
