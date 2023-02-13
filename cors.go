package httpbuilder

import (
	"fmt"
	"net/http"
)

type headerList []string

var wildcard = headerList{"*"}

func (l headerList) ToHeader() string {
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	default:
		return fmt.Sprintf("%s, %s", l[0], l[1:].ToHeader())
	}
}

type CorsConfig struct {
	AllowOrigins     headerList `yaml:"allowOrigins"`
	AllowHeaders     headerList `yaml:"allowHeaders"`
	AllowMethods     headerList `yaml:"allowMethods"`
	AllowCredentials headerList `yaml:"allowCredentials"`
}

type CorsConfigCallback func(*CorsConfig)

func Cors(userConfig CorsConfigCallback) Middleware {
	config := CorsConfig{
		AllowOrigins:     wildcard,
		AllowHeaders:     wildcard,
		AllowMethods:     wildcard,
		AllowCredentials: wildcard}
	if userConfig != nil {
		userConfig(&config)
	}
	return ToMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.AllowOrigins.ToHeader())
		w.Header().Set("Access-Control-Allow-Headers", config.AllowHeaders.ToHeader())
		w.Header().Set("Access-Control-Allow-Methods", config.AllowMethods.ToHeader())
		w.Header().Set("Access-Control-Allow-Credentials", config.AllowCredentials.ToHeader())
	})
}
