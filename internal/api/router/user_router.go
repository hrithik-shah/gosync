package router

import "gosync/internal/api/controller"

func UserRoutes(r *AppRouter, ctrl *controller.UserController) {
	r.Get("/me", ctrl.GetProfile)
	r.Patch("/me", ctrl.UpdateProfile)
}
