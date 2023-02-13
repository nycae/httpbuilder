package httpbuilder

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Builder struct {
	middlewares []Middleware
	handler     http.HandlerFunc
}

func (b *Builder) WithMiddleware(m Middleware) *Builder {
	b.middlewares = append(b.middlewares, m)
	return b
}

func (b *Builder) Build() http.HandlerFunc {
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		b.handler = b.middlewares[i](b.handler)
	}
	return b.handler
}

func FromFunc(fn http.HandlerFunc) *Builder {
	return &Builder{
		handler:     fn,
		middlewares: []Middleware{},
	}
}

func FromHandler(handler http.Handler) *Builder {
	return &Builder{
		handler:     handler.ServeHTTP,
		middlewares: []Middleware{},
	}
}

func ToMiddleware(handler http.HandlerFunc) Middleware {
	return func(innerHandler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
			innerHandler(w, r)
		}
	}
}
