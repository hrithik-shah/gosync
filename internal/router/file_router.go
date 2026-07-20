package router

import "gosync/internal/controller"

func FileRoutes(r *AppRouter, fileCtrl *controller.FileController) {
	r.Post("/", fileCtrl.Upload)
	r.Get("/{id}", fileCtrl.Download)
	r.Delete("/{id}", fileCtrl.Delete)
}
