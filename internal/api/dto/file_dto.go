package dto

import "gosync/internal/models"

type FileDTO struct {
	ID          string
	Name        string
	DirectoryID string
	Hash        string
}

func ToFileDTO(f *models.File) FileDTO {
	return FileDTO{
		ID:          f.ID.String(),
		Name:        f.Name,
		DirectoryID: f.DirectoryID.String(),
		Hash:        f.ContentHash,
	}
}
