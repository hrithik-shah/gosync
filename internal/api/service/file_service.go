package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gosync/internal/api/apperror"
	"gosync/internal/api/dto"
	"gosync/internal/api/repository"
	"gosync/internal/models"
	"gosync/internal/storage"
)

type FileService struct {
	q       *repository.Query
	storage storage.Storage
}

func NewFileService(q *repository.Query, store storage.Storage) *FileService {
	return &FileService{q: q, storage: store}
}

// Upload creates a new file record with an initial version, storing the
// content bytes via the storage backend.
func (s *FileService) Upload(ctx context.Context, userID, directoryID uuid.UUID, name string, content io.Reader) (*models.File, error) {
	// Ownership check on the target directory before creating anything.
	_, err := s.q.Directory.WithContext(ctx).
		Where(s.q.Directory.ID.Eq(directoryID.String()), s.q.Directory.UserID.Eq(userID.String())).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("directory not found")
		}
		return nil, apperror.Internal("failed to look up directory", err)
	}

	fileID := uuid.New()
	versionID := uuid.New()
	storageKey := fmt.Sprintf("%s/%s", fileID, versionID)

	hasher := sha256.New()
	size, err := s.storage.Save(storageKey, io.TeeReader(content, hasher))
	if err != nil {
		return nil, apperror.Internal("failed to store file content", err)
	}
	contentHash := hex.EncodeToString(hasher.Sum(nil))

	file := &models.File{
		ID:          fileID,
		UserID:      userID,
		DirectoryID: directoryID,
		Name:        name,
		ContentHash: contentHash,
	}

	err = s.q.Transaction(func(tx *repository.Query) error {
		txCtx := tx.WithContext(ctx)

		if err := txCtx.File.Create(file); err != nil {
			return err
		}

		version := &models.FileVersion{
			ID:            versionID,
			FileID:        fileID,
			VersionNumber: 1,
			StorageKey:    storageKey,
			SizeBytes:     size,
			ContentHash:   contentHash,
		}
		if err := txCtx.FileVersion.Create(version); err != nil {
			return err
		}

		_, err := txCtx.File.
			Where(tx.File.ID.Eq(fileID.String())).
			Update(tx.File.CurrentVersionID, versionID.String())
		return err
	})
	if err != nil {
		// best-effort cleanup of the orphaned blob if the DB transaction failed
		_ = s.storage.Delete(storageKey)
		return nil, apperror.Internal("failed to create file", err)
	}

	file.CurrentVersionID = &versionID
	return file, nil
}

// GetMetadata returns a file's metadata, verifying ownership.
func (s *FileService) GetMetadata(ctx context.Context, userID, fileID uuid.UUID) (dto.FileDTO, error) {
	file, err := s.getOwnedFile(ctx, userID, fileID)
	if err != nil {
		return dto.FileDTO{}, err
	}
	return dto.FileDTO{
		ID:          file.ID.String(),
		DirectoryID: file.DirectoryID.String(),
		Name:        file.Name,
		Hash:        file.ContentHash,
	}, nil
}

// Update renames a file.
func (s *FileService) Update(ctx context.Context, userID, fileID uuid.UUID, name string) (*models.File, error) {
	file, err := s.getOwnedFile(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	_, err = s.q.File.WithContext(ctx).
		Where(s.q.File.ID.Eq(fileID.String())).
		Update(s.q.File.Name, name)
	if err != nil {
		return nil, apperror.Internal("failed to update file", err)
	}

	file.Name = name
	return file, nil
}

// Move changes a file's parent directory, verifying ownership of both
// the file and the destination directory.
func (s *FileService) Move(ctx context.Context, userID, fileID, newDirectoryID uuid.UUID) (*models.File, error) {
	file, err := s.getOwnedFile(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	_, err = s.q.Directory.WithContext(ctx).
		Where(s.q.Directory.ID.Eq(newDirectoryID.String()), s.q.Directory.UserID.Eq(userID.String())).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("destination directory not found")
		}
		return nil, apperror.Internal("failed to look up destination directory", err)
	}

	_, err = s.q.File.WithContext(ctx).
		Where(s.q.File.ID.Eq(fileID.String())).
		Update(s.q.File.DirectoryID, newDirectoryID.String())
	if err != nil {
		return nil, apperror.Internal("failed to move file", err)
	}

	file.DirectoryID = newDirectoryID
	return file, nil
}

// Delete soft-deletes a file (sets DeletedAt via GORM's soft delete).
func (s *FileService) Delete(ctx context.Context, userID, fileID uuid.UUID) error {
	file, err := s.getOwnedFile(ctx, userID, fileID)
	if err != nil {
		return err
	}

	if _, err := s.q.File.WithContext(ctx).Delete(file); err != nil {
		return apperror.Internal("failed to delete file", err)
	}
	return nil
}

// UploadNewVersion adds a new version to an existing file and promotes
// it to the current version.
func (s *FileService) UploadNewVersion(ctx context.Context, userID, fileID uuid.UUID, content io.Reader) (*models.FileVersion, error) {
	if _, err := s.getOwnedFile(ctx, userID, fileID); err != nil {
		return nil, err
	}

	latest, err := s.q.FileVersion.WithContext(ctx).
		Where(s.q.FileVersion.FileID.Eq(fileID.String())).
		Order(s.q.FileVersion.VersionNumber.Desc()).
		First()
	if err != nil {
		return nil, apperror.Internal("failed to look up latest version", err)
	}

	versionID := uuid.New()
	storageKey := fmt.Sprintf("%s/%s", fileID, versionID)

	hasher := sha256.New()
	size, err := s.storage.Save(storageKey, io.TeeReader(content, hasher))
	if err != nil {
		return nil, apperror.Internal("failed to store file content", err)
	}
	contentHash := hex.EncodeToString(hasher.Sum(nil))

	version := &models.FileVersion{
		ID:            versionID,
		FileID:        fileID,
		VersionNumber: latest.VersionNumber + 1,
		StorageKey:    storageKey,
		SizeBytes:     size,
		ContentHash:   contentHash,
	}

	err = s.q.Transaction(func(tx *repository.Query) error {
		txCtx := tx.WithContext(ctx)

		if err := txCtx.FileVersion.Create(version); err != nil {
			return err
		}

		_, err := txCtx.File.
			Where(tx.File.ID.Eq(fileID.String())).
			Updates(map[string]any{
				"current_version_id": versionID.String(),
				"content_hash":       contentHash,
			})
		return err
	})
	if err != nil {
		_ = s.storage.Delete(storageKey)
		return nil, apperror.Internal("failed to save new version", err)
	}

	return version, nil
}

// GetContent returns a reader for a file's current version content.
func (s *FileService) GetContent(ctx context.Context, userID, fileID uuid.UUID) (io.ReadCloser, error) {
	file, err := s.getOwnedFile(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}
	if file.CurrentVersionID == nil {
		return nil, apperror.NotFound("file has no content")
	}

	version, err := s.q.FileVersion.WithContext(ctx).
		Where(s.q.FileVersion.ID.Eq(file.CurrentVersionID.String())).
		First()
	if err != nil {
		return nil, apperror.Internal("failed to look up current version", err)
	}

	reader, err := s.storage.Open(version.StorageKey)
	if err != nil {
		return nil, apperror.Internal("failed to open file content", err)
	}
	return reader, nil
}

// GetVersionContent returns a reader for a specific version's content.
func (s *FileService) GetVersionContent(ctx context.Context, userID, fileID uuid.UUID, versionNumber int) (io.ReadCloser, error) {
	if _, err := s.getOwnedFile(ctx, userID, fileID); err != nil {
		return nil, err
	}

	version, err := s.q.FileVersion.WithContext(ctx).
		Where(
			s.q.FileVersion.FileID.Eq(fileID.String()),
			s.q.FileVersion.VersionNumber.Eq(int32(versionNumber)),
		).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("file version not found")
		}
		return nil, apperror.Internal("failed to look up file version", err)
	}

	reader, err := s.storage.Open(version.StorageKey)
	if err != nil {
		return nil, apperror.Internal("failed to open file content", err)
	}
	return reader, nil
}

// ListVersions returns all versions of a file, verifying ownership.
func (s *FileService) ListVersions(ctx context.Context, userID, fileID uuid.UUID) ([]*models.FileVersion, error) {
	if _, err := s.getOwnedFile(ctx, userID, fileID); err != nil {
		return nil, err
	}

	versions, err := s.q.FileVersion.WithContext(ctx).
		Where(s.q.FileVersion.FileID.Eq(fileID.String())).
		Order(s.q.FileVersion.VersionNumber).
		Find()
	if err != nil {
		return nil, apperror.Internal("failed to list file versions", err)
	}
	return versions, nil
}

// getOwnedFile loads a file and verifies it belongs to userID, returning
// a NotFound error either if it doesn't exist or belongs to someone else
// (never reveal existence of another user's file).
func (s *FileService) getOwnedFile(ctx context.Context, userID, fileID uuid.UUID) (*models.File, error) {
	file, err := s.q.File.WithContext(ctx).
		Where(s.q.File.ID.Eq(fileID.String()), s.q.File.UserID.Eq(userID.String())).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("file not found")
		}
		return nil, apperror.Internal("failed to look up file", err)
	}
	return file, nil
}
