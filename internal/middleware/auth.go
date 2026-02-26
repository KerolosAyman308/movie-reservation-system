package middleware

import (
	"context"
	"net/http"
	"strings"

	"movie/system/internal/auth"
	"movie/system/internal/user"
	dtos "movie/system/internal/user/DTOs"
	"movie/system/pkg"
)

type contextKey string

// UserCtxKey is the key used to store the authenticated user in the request context.
const UserCtxKey contextKey = "userKey"

// AuthMiddleware holds dependencies for the JWT authentication middleware.
type AuthMiddleware struct {
	Authenticator auth.Authenticator
	UserService   user.Service
}

// NewAuthMiddleware creates an AuthMiddleware.
func NewAuthMiddleware(authenticator auth.Authenticator, service user.Service) *AuthMiddleware {
	return &AuthMiddleware{
		Authenticator: authenticator,
		UserService:   service,
	}
}

// Authenticate validates the Bearer JWT token and injects the user into the context.
// Routes that use this middleware can retrieve the user with UserFromContext.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkg.Unauthorized(w, r, (*any)(nil))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkg.Unauthorized(w, r, (*any)(nil))
			return
		}

		jwtToken, err := m.Authenticator.ValidateAccessToken(parts[1])
		if err != nil || !jwtToken.Valid {
			pkg.Unauthorized(w, r, (*any)(nil))
			return
		}

		// Use the typed CustomClaims — no fragile MapClaims string-formatting needed.
		claims, ok := jwtToken.Claims.(*auth.CustomClaims)
		if !ok {
			pkg.Unauthorized(w, r, (*any)(nil))
			return
		}

		userDto, err := m.UserService.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			pkg.Unauthorized(w, r, (*any)(nil))
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, userDto)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin checks that the authenticated user is an admin.
// Must be chained after Authenticate.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFromContext(r.Context())
		if !ok || !u.IsAdmin {
			pkg.Forbidden(w, r, (*any)(nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserFromContext retrieves the authenticated user DTO from the context.
func UserFromContext(ctx context.Context) (*dtos.UserResponse, bool) {
	u, ok := ctx.Value(UserCtxKey).(*dtos.UserResponse)
	return u, ok
}
