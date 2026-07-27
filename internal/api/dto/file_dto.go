package dto

type UpdateFileRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type MoveFileRequest struct {
	NewDirectoryID string `json:"directory_id" validate:"required,uuid"`
}
