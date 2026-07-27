package router

import "gosync/internal/api/controller"

func DeviceRoutes(r *AppRouter, deviceCtrl *controller.DeviceController) {
	r.Post("/", deviceCtrl.Create)
	r.Get("/", deviceCtrl.List)
}
