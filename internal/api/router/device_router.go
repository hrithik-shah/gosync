package router

import "gosync/internal/controller"

func DeviceRoutes(r *AppRouter, deviceCtrl *controller.DeviceController) {
	r.Post("/", deviceCtrl.Create)
	r.Get("/", deviceCtrl.List)
}
