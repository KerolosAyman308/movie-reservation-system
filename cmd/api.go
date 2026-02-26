package main

import (
	"movie/system/internal/auth"
	"movie/system/internal/config"
	"movie/system/internal/middleware"
	"movie/system/internal/user"
	"movie/system/pkg"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	pkg.Ok("Hello world", "", w)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	pkg.NotFound(w, r, (*any)(nil))
}

// InitializeAPI wires all dependencies and returns the configured router.
func InitializeAPI(db *gorm.DB, cfg config.Config) *chi.Mux {
	router := chi.NewRouter()

	// Global middlewares
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	// Wire dependencies
	authenticator := auth.NewJWTAuthenticator(cfg)
	repo := user.NewRepository(db)
	userService := user.NewService(repo)

	authMw := middleware.NewAuthMiddleware(authenticator, userService)
	userHandler := user.NewUserHandler(userService)
	authHandler := auth.NewAuthHandler(userService, authenticator, cfg.IsProduction)

	router.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/", handleMain)

		// Auth routes — public
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/signup", authHandler.SignUp)
			r.Post("/refresh", authHandler.RefreshToken)
		})

		// User routes — require authentication
		r.Route("/users", func(r chi.Router) {
			r.Use(authMw.Authenticate)

			r.With(middleware.RequireAdmin).Get("/", userHandler.GetAllUsers)
			r.With(middleware.RequireAdmin).Post("/", userHandler.AddUser)
			r.Get("/{userId}", userHandler.GetUserById)
			r.With(middleware.RequireAdmin).Patch("/{userId}/role", userHandler.ChangeRole)
		})
	})
	router.HandleFunc("/*", handleNotFound)

	return router
}
