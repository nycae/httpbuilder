package httpbuilder

import (
	"net/http"
)

type CorsConfig struct {
	AllowOrigins string `yaml:"allowOrigins"`
}

type CorsConfigCallback func(*CorsConfig)

func Cors(userConfig CorsConfigCallback) Middleware {
	config := CorsConfig{AllowOrigins: "*"}
	if userConfig != nil {
		userConfig(&config)
	}
	return ToMiddleware(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Access-Control-Allow-Origin", config.AllowOrigins)
	})
}
