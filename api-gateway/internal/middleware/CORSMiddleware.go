package middleware

import (
	"log"
	"net/http"
)

// CORSOptions определяет разрешенные источники CORS
var CORSOptions = map[string]bool{
	"http://localhost:5173": true, // Разрешенный источник
	// Добавьте сюда другие разрешенные источники
}

// CORSMiddleware добавляет необходимые заголовки CORS
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("CORS Request from Origin: %s", origin) // Логируем источник

		// Проверяем, разрешен ли источник
		if _, ok := CORSOptions[origin]; ok {
			rw.Header().Set("Access-Control-Allow-Origin", origin) // Устанавливаем разрешенный источник
		} else {
			rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Устанавливаем запасной источник или блокируем
		}

		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обрабатываем preflight запросы
		if r.Method == http.MethodOptions {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		// Передаем управление следующему обработчику
		next.ServeHTTP(rw, r)
	}
}
