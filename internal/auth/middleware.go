package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/open-tfe/tfe-service/internal/constants"
	"go.uber.org/zap"
)

func JWTMiddleware(secret string, logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("Processing request", zap.String("path", r.URL.Path))

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Debug("Missing Authorization header")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			logger.Debug("Found Authorization header, parsing token")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				logger.Debug("Token validation failed", zap.Error(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			logger.Debug("Token is valid, checking claims")
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logger.Debug("Failed to parse token claims")
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			email, ok := claims["email"].(string)
			if !ok || email == "" {
				logger.Debug("No valid email found in token claims")
				http.Error(w, "Invalid email in token", http.StatusUnauthorized)
				return
			}

			logger.Debug("Successfully authenticated user", zap.String("email", email))
			ctx := context.WithValue(r.Context(), constants.UserEmailKey, email)
			ctx = context.WithValue(ctx, constants.UserTokenKey, tokenString)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
