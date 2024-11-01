package jwt

import (
	userModel "SpotifySorter/models"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type User interface {
	GetUserByAccessToken(accessToken string) (*userModel.User, error)
}

const UserContextKey = "user"

func JWTMiddleware(secret string, user User) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Токен должен быть в формате Bearer <token>
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]

			// Проверяем валидность токена
			claims := &userModel.Claims{}
			_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Проверяем алгоритм подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorMalformed)
				}
				return []byte(secret), nil
			})

			if err != nil {
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorExpired != 0 {
						http.Error(w, "token expired", http.StatusUnauthorized)
						return
					}
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем данные о пользователе в контекст запроса
			user, err := user.GetUserByAccessToken(tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Set user in request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) *userModel.User {
	user, _ := ctx.Value(UserContextKey).(*userModel.User)
	return user
}
