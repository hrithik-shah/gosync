package payload

import "gosync/internal/api/dto"

type FileInfo struct {
	ID          string `json:"id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required,min=1"`
	DirectoryID string `json:"directory_id" validate:"required,uuid"`
	Hash        string `json:"hash" validate:"required,len=256"`
}

type UpdateFileRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type MoveFileRequest struct {
	NewParentDirectoryID string `json:"directory_id" validate:"required,uuid"`
}

func FromFileDTO(file dto.FileDTO) FileInfo {
	return FileInfo{ID: file.ID, Name: file.Name, DirectoryID: file.DirectoryID, Hash: file.Hash}
}

func FromFileDTOSlice(files []dto.FileDTO) []FileInfo {
	result := make([]FileInfo, len(files))

	for i, file := range files {
		result[i] = FromFileDTO(file)
	}

	return result
}
