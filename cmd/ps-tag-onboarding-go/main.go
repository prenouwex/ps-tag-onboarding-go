package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/wexinc/ps-tag-onboarding-go/internal/constants"
	"github.com/wexinc/ps-tag-onboarding-go/internal/controller"
	"github.com/wexinc/ps-tag-onboarding-go/internal/database"
	"github.com/wexinc/ps-tag-onboarding-go/internal/log"
	"github.com/wexinc/ps-tag-onboarding-go/internal/repository"
	"github.com/wexinc/ps-tag-onboarding-go/internal/router"
	"github.com/wexinc/ps-tag-onboarding-go/internal/service"
	"net/http"

	"time"
)

func main() {
	db := database.CreateNewGormDB()
	userRepository := repository.UserRepository{DB: db}
	userValidation := service.UserValidationService{&userRepository}
	userService := service.UserService{&userRepository, &userValidation}
	userController := controller.UserController{&userService}
	userRoutes := router.UserRoutes{&userController}
	handleRequests(&userRoutes)
}

func handleRequests(userRoutes *router.UserRoutes) {
	r := chi.NewRouter()

	// Config
	r.Use(middleware.RequestID)
	r.Use(log.RequestLogger)
	r.Use(log.RequestFileLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Content-Type", "JWT-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Routes
	//r.Group(protectedRoutes)
	r.Group(userRoutes.UserRoutes)

	//Run
	httpPort := fmt.Sprintf(":%s", constants.HTTP_PORT)

	log.Info.Printf("Starting server on %v\n", httpPort)
	log.Error.Println(http.ListenAndServe(httpPort, r))
}
