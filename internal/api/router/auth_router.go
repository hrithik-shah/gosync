package router

import (
	"gosync/internal/api/controller"
	"gosync/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r *AppRouter, authCtrl *controller.AuthController) {
	r.Post("/login", authCtrl.Login)
	r.Post("/register", authCtrl.Register)

	r.Router.Group(func(sub chi.Router) {
		sub.Use(middleware.RequireAuth)
		wrapped := &AppRouter{sub}
		wrapped.Post("/logout", authCtrl.Logout)
	})
}
