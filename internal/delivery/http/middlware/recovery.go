package middlware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery это middleware для восстановления после паники
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Логируем ошибку и стек вызовов
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Возвращаем HTTP 500
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := map[string]string{
					"error": "Внутренняя ошибка сервера",
				}

				// Игнорируем ошибку кодирования, так как это уже обработка ошибки
				_ = json.NewEncoder(w).Encode(response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
