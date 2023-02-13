package httpbuilder

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Cors(t *testing.T) {
	testCases := map[string]struct {
		endpoint http.HandlerFunc
		config   CorsConfigCallback
		expected map[string]string
	}{
		"all_wildcard": {
			endpoint: emptyHandlerFunc,
			config:   nil,
			expected: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Headers":     "*",
				"Access-Control-Allow-Methods":     "*",
				"Access-Control-Allow-Credentials": "*",
			},
		},
		"single_method": {
			endpoint: emptyHandlerFunc,
			config: func(c *CorsConfig) {
				c.AllowMethods = []string{http.MethodGet}
			},
			expected: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Headers":     "*",
				"Access-Control-Allow-Methods":     "GET",
				"Access-Control-Allow-Credentials": "*",
			},
		},
		"multiple_methods": {
			endpoint: emptyHandlerFunc,
			config: func(c *CorsConfig) {
				c.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut}
			},
			expected: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Headers":     "*",
				"Access-Control-Allow-Methods":     "GET, POST, PUT",
				"Access-Control-Allow-Credentials": "*",
			},
		},
		"all_access": {
			endpoint: emptyHandlerFunc,
			config: func(c *CorsConfig) {
				c.AllowMethods = headerList{http.MethodGet, http.MethodPost, http.MethodPut}
				c.AllowOrigins = headerList{"localhost.com"}
				c.AllowHeaders = headerList{"Content-Type", "Accept", "Connection", "Upgrade"}
				c.AllowCredentials = headerList{"true"}
			},
			expected: map[string]string{
				"Access-Control-Allow-Origin":      "localhost.com",
				"Access-Control-Allow-Headers":     "Content-Type, Accept, Connection, Upgrade",
				"Access-Control-Allow-Methods":     "GET, POST, PUT",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		"empty_header": {
			endpoint: emptyHandlerFunc,
			config: func(c *CorsConfig) {
				c.AllowMethods = []string{}
			},
			expected: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Headers":     "*",
				"Access-Control-Allow-Methods":     "",
				"Access-Control-Allow-Credentials": "*",
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
			endpoint := FromFunc(test.endpoint).
				WithMiddleware(Cors(test.config)).
				Build()
			endpoint(w, r)

			for key, expected := range test.expected {
				if got := w.Header().Get(key); got != expected {
					t.Errorf("Header Missmatch:\n\texpected: %v\n\tgot: %v", expected, got)
				}
			}
		})
	}
}
