package dto

type FileInfo struct {
	ID          string `json:"id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required,min=1"`
	DirectoryID string `json:"directory_id" validate:"required,uuid"`
}

type UpdateFileRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type MoveFileRequest struct {
	NewDirectoryID string `json:"directory_id" validate:"required,uuid"`
}
