package router

import "gosync/internal/controller"

func FileRoutes(r *AppRouter, fileCtrl *controller.FileController) {
	r.Post("/", fileCtrl.Upload)
	r.Get("/{id}", fileCtrl.GetMetadata)
	r.Patch("/{id}", fileCtrl.Update)
	r.Post("/{id}/move", fileCtrl.Move)
	r.Delete("/{id}", fileCtrl.Delete)

	r.Post("/{id}/content", fileCtrl.UploadNewFileContent)
	r.Get("/{id}/content", fileCtrl.GetFileContent)
	r.Get("/{id}/versions/{version}", fileCtrl.GetFileVersion)

	r.Get("/{id}/versions", fileCtrl.ListFileVersions)
	// r.Post("/{id}/versions/{version}/restore", fileCtrl.RestoreFileVersion) TODO
}
