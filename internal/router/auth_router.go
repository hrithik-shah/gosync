package router

import "gosync/internal/controller"

func AuthRoutes(r *AppRouter, authCtrl *controller.AuthController) {
	r.Post("/login", authCtrl.Login)
	r.Post("/register", authCtrl.Register)
	r.Post("/logout", authCtrl.Logout)
}
