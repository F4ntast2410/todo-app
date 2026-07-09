package middleware

import (
	"context"
	"net/http"
)

// Определяем кастомный тип для ключа контекста, чтобы избежать коллизий
type contextKey string

const UserIDKey contextKey = "userID"

// Интерфейс для взаимодействия мидлвари со слоем бизнес-логики
type AuthUsecase interface {
	GetUserIDBySession(ctx context.Context, sessionID string) (int, error)
}

func Auth(userUC AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Достаем куку сессии
			cookie, err := r.Cookie("session_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized) // Куки нет -> 401
				return
			}

			// 2. Проверяем сессию через UseCase в БД
			userID, err := userUC.GetUserIDBySession(r.Context(), cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized) // Сессия невалидна или истекла -> 401
				return
			}

			// 3. Кладем userID в контекст запроса и передаем запрос дальше
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
