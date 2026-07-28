package router

import "gosync/internal/api/controller"

func SyncRoutes(r *AppRouter, syncCtrl *controller.SyncController) {
	r.Post("/", syncCtrl.DetermineSyncActions)
	r.Get("/events", syncCtrl.GetEvents)
}
