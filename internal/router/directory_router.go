package router

import "gosync/internal/controller"

func DirectoryRoutes(r *AppRouter, dirCtrl *controller.DirectoryController) {
	r.Get("/", dirCtrl.List)
	r.Post("/", dirCtrl.Create)
	r.Delete("/{id}", dirCtrl.Delete)
}
