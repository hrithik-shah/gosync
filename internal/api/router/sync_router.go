package router

import "gosync/internal/controller"

func SyncRoutes(r *AppRouter, syncCtrl *controller.SyncController) {
	r.Post("/", syncCtrl.DetermineSyncActions)
	r.Get("/root", syncCtrl.GetRootHash)
	r.Post("/compare", syncCtrl.CompareHashes)
	r.Get("/tree/{directory_id}", syncCtrl.GetDirectoryTreeNode)
	r.Get("/events", syncCtrl.GetEvents)
}
