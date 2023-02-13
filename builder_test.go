package httpbuilder

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func emptyHandlerFunc(w http.ResponseWriter, r *http.Request) {}

func Test_FunctionBuilder(t *testing.T) {
	tests := map[string]struct {
		Endpoint    http.HandlerFunc
		Middlewares []Middleware
		Expected    string
	}{
		"happy path": {
			Endpoint: emptyHandlerFunc,
			Middlewares: []Middleware{
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("something"))
						handler(w, r)
					}
				},
			},
			Expected: "something",
		},
		"double middleware": {
			Endpoint: emptyHandlerFunc,
			Middlewares: []Middleware{
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("something"))
						handler(w, r)
					}
				},
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte(" else"))
						handler(w, r)
					}
				},
			},
			Expected: "something else",
		},
		"before and after": {
			Endpoint: emptyHandlerFunc,
			Middlewares: []Middleware{
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("something"))
						handler(w, r)
					}
				},
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						handler(w, r)
						w.Write([]byte(" else"))
					}
				},
			},
			Expected: "something else",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			builder := FromFunc(test.Endpoint)
			for _, middleware := range test.Middlewares {
				builder.WithMiddleware(middleware)
			}
			endpoint := builder.Build()
			w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
			endpoint(w, r)

			if received := string(w.Body.Bytes()); received != test.Expected {
				t.Errorf("error in test:\n\tExpected: %v\n\tReceived: %v", test.Expected, received)
			}
		})
	}
}

type emptyHandler struct{}

func (*emptyHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func Test_HandlerBuilder(t *testing.T) {
	tests := map[string]struct {
		Handler     http.Handler
		Middlewares []Middleware
		Expected    string
	}{
		"happy path": {
			Handler: &emptyHandler{},
			Middlewares: []Middleware{
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("something"))
						handler(w, r)
					}
				},
			},
			Expected: "something",
		},
		"double middleware": {
			Handler: &emptyHandler{},
			Middlewares: []Middleware{
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("something"))
						handler(w, r)
					}
				},
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte(" else"))
						handler(w, r)
					}
				},
			},
			Expected: "something else",
		},
		"before and after": {
			Handler: &emptyHandler{},
			Middlewares: []Middleware{
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("something"))
						handler(w, r)
					}
				},
				func(handler http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						handler(w, r)
						w.Write([]byte(" else"))
					}
				},
			},
			Expected: "something else",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			builder := FromHandler(test.Handler)
			for _, middleware := range test.Middlewares {
				builder.WithMiddleware(middleware)
			}
			endpoint := builder.Build()
			w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
			endpoint(w, r)

			if received := string(w.Body.Bytes()); received != test.Expected {
				t.Errorf("error in test:\n\tExpected: %v\n\tReceived: %v", test.Expected, received)
			}
		})
	}
}

func Test_ToMiddleware(t *testing.T) {
	tests := map[string]struct {
		Handler     http.Handler
		Middlewares []http.HandlerFunc
		Expected    string
	}{
		"happy path": {
			Handler: &emptyHandler{},
			Middlewares: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("something"))
				},
			},
			Expected: "something",
		},
		"double middleware": {
			Handler: &emptyHandler{},
			Middlewares: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("something "))
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("else"))
				},
			},
			Expected: "something else",
		},
		"before and after": {
			Handler: &emptyHandler{},
			Middlewares: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("something"))
				},
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(" else"))
				},
			},
			Expected: "something else",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			builder := FromHandler(test.Handler)
			for _, middleware := range test.Middlewares {
				builder.WithMiddleware(ToMiddleware(middleware))
			}
			endpoint := builder.Build()
			w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
			endpoint(w, r)

			if received := string(w.Body.Bytes()); received != test.Expected {
				t.Errorf("error in test:\n\tExpected: %v\n\tReceived: %v", test.Expected, received)
			}
		})
	}
}
