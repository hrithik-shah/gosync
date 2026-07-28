package dto

type SycnEventType string

const (
	Create SycnEventType = "create"
	Update SycnEventType = "update"
	Delete SycnEventType = "delete"
	Move   SycnEventType = "move"
)

type SyncEventDTO struct {
	ID          string
	Type        SycnEventType
	FileID      string
	DirectoryID string
}
