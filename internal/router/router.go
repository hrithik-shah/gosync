package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"gosync/internal/controller"
	"gosync/internal/middleware"
)

type AppRouter struct {
	chi.Router
}

func (ar *AppRouter) Get(pattern string, h middleware.HandlerFunc) {
	ar.Router.Get(pattern, middleware.ErrorHandler(h))
}
func (ar *AppRouter) Post(pattern string, h middleware.HandlerFunc) {
	ar.Router.Post(pattern, middleware.ErrorHandler(h))
}
func (ar *AppRouter) Put(pattern string, h middleware.HandlerFunc) {
	ar.Router.Put(pattern, middleware.ErrorHandler(h))
}
func (ar *AppRouter) Delete(pattern string, h middleware.HandlerFunc) {
	ar.Router.Delete(pattern, middleware.ErrorHandler(h))
}
func (ar *AppRouter) Patch(pattern string, h middleware.HandlerFunc) {
	ar.Router.Patch(pattern, middleware.ErrorHandler(h))
}

// Route overrides chi.Router.Route to automatically pass down an AppRouter instance
func (ar *AppRouter) Route(pattern string, fn func(sub *AppRouter)) {
	ar.Router.Route(pattern, func(sub chi.Router) {
		fn(&AppRouter{sub})
	})
}

func New(db *gorm.DB) http.Handler {
	baseRouter := chi.NewRouter()
	baseRouter.Use(chimw.Logger)
	baseRouter.Use(chimw.Recoverer)

	r := &AppRouter{baseRouter}

	r.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	userCtrl := controller.NewUserController(db)
	authCtrl := controller.NewAuthController(db)
	fileCtrl := controller.NewFileController(db)
	dirCtrl := controller.NewDirectoryController(db)

	r.Route("/auth", func(sub *AppRouter) {
		AuthRoutes(sub, authCtrl)
	})

	r.Router.Group(func(sub chi.Router) {
		sub.Use(middleware.RequireAuth)

		wrappedSub := &AppRouter{sub}
		wrappedSub.Route("/users", func(s *AppRouter) { UserRoutes(s, userCtrl) })
		wrappedSub.Route("/files", func(s *AppRouter) { FileRoutes(s, fileCtrl) })
		wrappedSub.Route("/directories", func(s *AppRouter) { DirectoryRoutes(s, dirCtrl) })
	})

	return r
}
