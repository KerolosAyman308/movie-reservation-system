package main

import (
	"context"
	"log"
	"movie/system/internal/auth"
	"movie/system/internal/config"
	f "movie/system/internal/files"
	"movie/system/internal/middleware"
	movies "movie/system/internal/movies/handlers"
	ms "movie/system/internal/movies/services"
	"movie/system/internal/user"
	"movie/system/pkg"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsc "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	pkg.Ok("Hello world", "", w)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	pkg.NotFound(w, r, (*any)(nil))
}

func configAws(myCfg config.Config, db *gorm.DB) (f.IFiles, error) {
	ctx := context.TODO()
	s3Config, err := awsc.LoadDefaultConfig(ctx,
		awsc.WithRegion("garage"),
		awsc.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(myCfg.File.AWSAccessKey, myCfg.File.AWSSecretKey, "")),
		awsc.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           myCfg.File.AWSHost, // Your Garage endpoint URL
					SigningRegion: "garage",
				}, nil
			})),
	)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	s3Client := s3.NewFromConfig(s3Config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(myCfg.File.AWSHost)
		o.Region = "garage"
		o.UsePathStyle = true
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})

	uploader := manager.NewUploader(s3Client)
	presignClient := s3.NewPresignClient(s3Client)
	fileService := f.NewS3Service(s3Client, db, uploader, presignClient)

	return fileService, nil
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

	//fileService := f.NewFileService(cfg, db)
	//awsService, _ := configAws(cfg, db)
	var storageService f.IFiles
	if cfg.File.UseFile {
		storageService = f.NewFileService(cfg, db)
	} else {
		storageService, _ = configAws(cfg, db)
	}

	genreService := ms.NewGenreService(db)
	movieService := ms.NewMoviesService(db, genreService, storageService, cfg.File.BucketName)

	authMw := middleware.NewAuthMiddleware(authenticator, userService)
	userHandler := user.NewUserHandler(userService)
	moviesHandler := movies.NewMovieHandler(genreService, movieService)
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

		r.Route("/genres", func(r chi.Router) {
			r.Get("/", moviesHandler.FindPaginated)
			r.With(authMw.Authenticate).With(middleware.RequireAdmin).Post("/", moviesHandler.CreateGenre)
			r.With(authMw.Authenticate).With(middleware.RequireAdmin).Delete("/{genreId}", moviesHandler.DeleteGenre)
		})

		r.Route("/movies", func(r chi.Router) {
			r.Get("/", moviesHandler.MoviePaginated)
			r.Group(func(r chi.Router) {
				r.Use(authMw.Authenticate, middleware.RequireAdmin)
				r.Post("/", moviesHandler.CreateMovie)
				r.Post("/{movieId}/image", moviesHandler.UploadImage)
				r.Post("/{movieId}/genres/add", moviesHandler.AddMovieGenres)
				r.Post("/{movieId}/genres/delete", moviesHandler.DeleteMovieGenres)
				r.Delete("/{movieId}", moviesHandler.DeleteMovie)
			})
		})
	})
	router.HandleFunc("/*", handleNotFound)

	return router
}
