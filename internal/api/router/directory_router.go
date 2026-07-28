package router

import "gosync/internal/api/controller"

func DirectoryRoutes(r *AppRouter, dirCtrl *controller.DirectoryController) {
	r.Post("/", dirCtrl.Create)
	r.Get("/root", dirCtrl.GetRootMetadata)
	r.Get("/{id}", dirCtrl.GetMetadata)
	r.Get("/{id}/list", dirCtrl.ListContents)
	r.Patch("/{id}", dirCtrl.Update)
	r.Post("/{id}/move", dirCtrl.Move)
	r.Delete("/{id}", dirCtrl.Delete)
}
