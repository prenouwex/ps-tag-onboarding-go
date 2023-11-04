package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/runtime/middleware"
	"github.com/wexinc/ps-tag-onboarding-go/internal/controller"
	"net/http"
	//"github.com/go-chi/jwtauth/v5"
)

type UserRoutes struct {
	Controller *controller.UserController
}

func (ur *UserRoutes) UserRoutes(r chi.Router) {

	r.Route("/users", func(r chi.Router) {
		r.With(paginate).Get("/", ur.Controller.ListUsers)
		r.Post("/", ur.Controller.SaveUser) // POST /users
		////r.Get("/search", SearchUsers) // GET /users/search

		r.Route("/{userId}", func(r chi.Router) {
			//r.Use(UserCtx)            // Load the *User on the request context
			r.Get("/", ur.Controller.GetUser)       // GET /users/123
			r.Put("/", ur.Controller.UpdateUser)    // PUT /users/123
			r.Delete("/", ur.Controller.DeleteUser) // DELETE /users/123
		})

	})

}

func (ur *UserRoutes) SwaggerRoutes(r chi.Router) {
	// Serve the Swagger JSON
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./api")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "swagger.yaml", Path: "doc"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/doc", sh1)
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}
