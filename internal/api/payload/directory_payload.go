package payload

type DirectoryInfo struct {
	ID                string `json:"id" validate:"required,uuid"`
	Name              string `json:"name" validate:"required,min=1"`
	ParentDirectoryID string `json:"parent_directory_id" validate:"required,uuid"`
}

type CreateDirectoryRequest struct {
	Name              string  `json:"name" validate:"required,min=1"`
	ParentDirectoryID *string `json:"parent_directory_id" validate:"required,uuid"`
}

type UpdateDirectoryRequest struct {
	NewName string `json:"name" validate:"required,min=1"`
}

type MoveDirectoryRequest struct {
	NewParentDirectoryID string `json:"parent_directory_id" validate:"required,uuid"`
}

type CreateDirectoryResponse struct {
	DirectoryID string `json:"id"`
}

type GetMetadataResponse = DirectoryInfo

type ListDirectoryContentsResponse struct {
	Files []FileInfo `json:"files" validate:"required,dive"`
}
