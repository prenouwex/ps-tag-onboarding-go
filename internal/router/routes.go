package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/wexinc/ps-tag-onboarding-go/internal/controller"
	"net/http"
	//"github.com/go-chi/jwtauth/v5"
)

type UserRoutes struct {
	//Service *service.UserService
	Service *controller.UserController
}

func (ur *UserRoutes) UserRoutes(r chi.Router) {

	r.Route("/users", func(r chi.Router) {
		r.With(paginate).Get("/", ur.Service.ListUsers)
		r.Post("/", ur.Service.SaveUser) // POST /users
		////r.Get("/search", SearchUsers) // GET /users/search

		r.Route("/{userId}", func(r chi.Router) {
			//r.Use(ArticleCtx)            // Load the *User on the request context
			r.Get("/", ur.Service.GetUser)       // GET /users/123
			r.Put("/", ur.Service.UpdateUser)    // PUT /users/123
			r.Delete("/", ur.Service.DeleteUser) // DELETE /users/123
		})

	})

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
