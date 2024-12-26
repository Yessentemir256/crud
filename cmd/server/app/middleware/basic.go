package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

// Basic возвращает Middleware для проверки Basic-аутентификации.
func Basic(authFunc func(ctx context.Context, login, password string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Декодируем базовый логин:пароль
			encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
			decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
			if err != nil {
				http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
				return
			}

			// Разбиваем на логин и пароль
			credentials := strings.SplitN(string(decodedBytes), ":", 2)
			if len(credentials) != 2 {
				http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
				return
			}

			login, password := credentials[0], credentials[1]

			// Вызываем функцию аутентификации
			if !authFunc(r.Context(), login, password) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Если всё хорошо, передаем управление следующему обработчику
			next.ServeHTTP(w, r)
		})
	}
}
