package router

import "gosync/internal/controller"

func UserRoutes(r *AppRouter, ctrl *controller.UserController) {
	r.Get("/", ctrl.GetProfile)
	r.Post("/update", ctrl.UpdateProfile)
}
