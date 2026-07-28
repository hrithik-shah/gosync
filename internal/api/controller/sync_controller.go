package controller

import (
	"encoding/json"
	"gosync/internal/api/apperror"
	"gosync/internal/api/middleware"
	"gosync/internal/api/payload"
	"gosync/internal/api/utils/httputil"
	"net/http"
)

type SyncController struct {
	syncService *service.SyncService
}

func NewSyncController() *SyncController {
	return &SyncController{syncService: service.NewSyncService()}
}

// DetermineSyncActions godoc
// @Description  Determines the sync actions needed for a file
// @Summary      Determine sync actions
// @Tags         sync
// @Accept       json
// @Param        body  body      payload.DetermineSyncActionsRequest  true  "File details"
// @Success      200  {object}  payload.DetermineSyncActionsResponse
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /sync [post]
func (c *SyncController) DetermineSyncActions(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	var req payload.DetermineSyncActionsRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	action, err := c.syncService.DetermineSyncActions(userID, req.LastSyncEventID, req.RootHash)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(payload.DetermineSyncActionsResponse{Type: action})
}

// GetEvents godoc
// @Description  Gets events for a user
// @Summary      Get events
// @Tags         sync
// @Produce      json
// @Param        cursor  query     string  true  "Last-seen event ID by client"
// @Param        limit   query     int     false  "Max events to return (default 10, max 15)"
// @Success      200  {object}  payload.GetEventsResponse
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /sync/events [get]
func (c *SyncController) GetEvents(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	cursor, limit, err := httputil.ParsePagination(r)
	if err != nil {
		return err
	}

	events, count, err := c.syncService.GetEvents(userID, cursor, limit)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(payload.GetEventsResponse{
		Events: payload.FromSyncEventDTOSlice(events),
		Count:  count,
	})
}
