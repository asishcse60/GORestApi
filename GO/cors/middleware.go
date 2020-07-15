package cors

import (
	"net/http"
)

func Middleware(handler http.Handler)  http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(writer, request)
	})
}

