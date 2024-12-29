package middleware

import (
	"log"
	"net/http"
)

func Logger(handler http.Handler) http.Handler {
	// middleware "заворачивает" Handler, возвращая новый Handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// до выполнения handler'а
		log.Printf("START: %s %s", request.Method, request.URL.Path)

		// выполнение handler'а
		handler.ServeHTTP(writer, request)

		// после выполнения handler'а
		log.Printf("FINISH: %s %s", request.Method, request.URL.Path)
	})
}
