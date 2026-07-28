package dto

type EventType string

const (
	Create EventType = "create"
	Update EventType = "update"
	Delete EventType = "delete"
	Move   EventType = "move"
)

type EventDTO struct {
	ID          string
	Type        EventType
	FileID      string
	DirectoryID string
}
