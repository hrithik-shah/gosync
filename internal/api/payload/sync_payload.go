package payload

import "gosync/internal/api/dto"

type SyncActionType string

const (
	Event  SyncActionType = "event"
	Merkle SyncActionType = "merkle"
)

type SycnEventType string

const (
	Create SycnEventType = "create"
	Update SycnEventType = "update"
	Delete SycnEventType = "delete"
	Move   SycnEventType = "move"
)

type SyncEventInfo struct {
	ID          string `json:"id" validate:"required"`
	Type        string `json:"type" validate:"required"`
	FileID      string `json:"file_id" validate:"omitempty,uuid"`
	DirectoryID string `json:"directory_id" validate:"omitempty,uuid"`
}

type DetermineSyncActionsRequest struct {
	LastSyncEventID string `json:"last_sync_event_id" validate:"required"`
	RootHash        string `json:"root_hash" validate:"required"`
}

type DetermineSyncActionsResponse struct {
	Type SyncActionType `json:"type" validate:"required"`
}

type GetEventsResponse struct {
	Events []SyncEventInfo `json:"events" validate:"required,dive"`
	Count  int             `json:"count" validate:"required"`
}

func FromSyncEventDTO(eventDTO dto.SyncEventDTO) SyncEventInfo {
	return SyncEventInfo{
		ID:          eventDTO.ID,
		Type:        string(eventDTO.Type),
		FileID:      eventDTO.FileID,
		DirectoryID: eventDTO.DirectoryID,
	}
}

func FromSyncEventDTOSlice(eventDTOs []dto.SyncEventDTO) []SyncEventInfo {
	events := make([]SyncEventInfo, len(eventDTOs))
	for i, eventDTO := range eventDTOs {
		events[i] = FromSyncEventDTO(eventDTO)
	}
	return events
}
