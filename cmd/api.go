package main

import (
	user "movie/system/internal/user"
	"movie/system/pkg"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type IHandler struct {
	DB *gorm.DB
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	pkg.Ok("Hello world", "", w)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	pkg.NotFound(w, r, (*any)(nil))
}

func initUserHandlers(router *chi.Mux, db *gorm.DB) {
	var handle user.UserHandler = user.UserHandler{
		UserService: user.UserService{
			DB: db,
		},
	}

	router.Route("/users", func(r chi.Router) {
		r.Get("/", handle.GetAllUsers)
		r.Post("/", handle.AddUser)

		r.Get("/{userId}", handle.GetUserById)
		r.Post("/{userId}", handle.ChangeRole)
	})
}

func initMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}

func InitializeAPI(db *gorm.DB) *chi.Mux {
	router := chi.NewRouter()

	initMiddlewares(router)
	initUserHandlers(router, db)

	router.HandleFunc("/", handleMain)
	router.HandleFunc("/*", handleNotFound)

	return router
}
