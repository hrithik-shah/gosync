package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"gosync/internal/api/apperror"
	"gosync/internal/api/middleware"
	"gosync/internal/api/payload"
	"gosync/internal/api/service"
	"gosync/internal/api/utils/httputil"

	"gosync/internal/config"
)

type FileController struct {
	fileService *service.FileService
}

func NewFileController() *FileController {
	return &FileController{fileService: service.NewFileService()}
}

// Upload godoc
// @Description  Uploads a new file and associates it with a directory
// @Summary      Upload a new file
// @Tags         files
// @Accept       multipart/form-data
// @Param        directory_id  formData  string  true  "Target directory ID"
// @Param        file          formData  file    true  "File content"
// @Success      201  {object}  models.File
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files [post]
func (c *FileController) Upload(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	if err := r.ParseMultipartForm(config.MaxFileMemory); err != nil {
		return apperror.BadRequest("invalid multipart form")
	}

	directoryID, err := httputil.ParseUUIDParam(r, "directory_id")
	if err != nil {
		return apperror.BadRequest("invalid or missing directory_id")
	}

	f, header, err := r.FormFile("file")
	if err != nil {
		return apperror.BadRequest("missing file")
	}
	defer f.Close()

	fileDTO, err := c.fileService.Upload(r.Context(), userID, directoryID, header.Filename, f)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(payload.FromFileDTO(fileDTO))
}

// GetMetadata godoc
// @Description  Get file metadata
// @Summary      Returns metadata for a file owned by the authenticated user
// @Tags         files
// @Produce      json
// @Param        id  path  string  true  "File ID"
// @Success      200  {object}  models.File
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id} [get]
func (c *FileController) GetMetadata(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	metadata, err := c.fileService.GetMetadata(userID, fileID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(metadata)
}

// Update handles PATCH /{id} — renames a file.
// Update godoc
// @Description  Rename a file
// @Summary      Updates name metadata for a file owned by the authenticated user
// @Tags         files
// @Accept       json
// @Param        id    	path  string             true  "File ID"
// @Param        body  	body  payload.UpdateFileRequest  true  "New file name"
// @Success      200
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id} [patch]
func (c *FileController) Update(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	var req payload.UpdateFileRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	if err := c.fileService.Update(userID, fileID, req.Name); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Move godoc
// @Description  Move a file to a different directory
// @Summary      Moves a file to a different directory. Both file and target directory must be owned by the authenticated user.
// @Tags         files
// @Accept       json
// @Param        id    	path  string             true  "File ID"
// @Param        body  	body  payload.MoveFileRequest  true  "New directory ID"
// @Success      200
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id}/move [post]
func (c *FileController) Move(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	var req payload.MoveFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}

	newDirID, err := uuid.Parse(req.NewParentDirectoryID)
	if err != nil {
		return apperror.BadRequest("invalid directory_id")
	}

	err = c.fileService.Move(userID, fileID, newDirID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Delete godoc
// @Description  Delete a file
// @Summary      Deletes a file owned by the authenticated user
// @Tags         files
// @Param        id    	path  string             true  "File ID"
// @Success      200
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id} [delete]
func (c *FileController) Delete(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	if err := c.fileService.Delete(userID, fileID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// UploadNewFileContent godoc
// @Summary      Uploads new file content
// @Description  Uploads new file content, and keeps it associated with the same directory.
// @Tags         files
// @Accept       multipart/form-data
// @Param        id    path      string  true  "Target file ID"
// @Param        file  formData  file    true  "New file content"
// @Success      201  {object}  map[string]int  "version_number"
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id}/content [post]
func (c *FileController) UploadNewFileContent(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	if err := r.ParseMultipartForm(config.MaxFileMemory); err != nil {
		return apperror.BadRequest("invalid multipart form")
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		return apperror.BadRequest("missing file")
	}
	defer f.Close()

	versionNumber, err := c.fileService.UploadNewVersion(userID, fileID, f)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(map[string]int{
		"version_number": versionNumber,
	})
}

// GetFileContent godoc
// @Summary      Gets file content
// @Description  Streams the current version of a file.
// @Tags         files
// @Produce      octet-stream
// @Param        id    path      string  true  "Target file ID"
// @Success      200  {string}  binary  "Raw file content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id}/content [get]
func (c *FileController) GetFileContent(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	reader, err := c.fileService.GetContent(userID, fileID)
	if err != nil {
		return err
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	_, copyErr := io.Copy(w, reader)
	return copyErr
}

// GetFileVersion godoc
// @Summary      Gets file content for a specific version
// @Description  Streams the content of a specific version of a file.
// @Tags         files
// @Produce      octet-stream
// @Param        id      path      string  true  "Target file ID"
// @Param        version path      int     true  "Version number"
// @Success      200     {string}  binary  "Raw file content"
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id}/versions/{version} [get]
func (c *FileController) GetFileVersion(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	versionNumber, err := strconv.Atoi(chi.URLParam(r, "version"))
	if err != nil {
		return apperror.BadRequest("invalid version number")
	}

	reader, err := c.fileService.GetVersionContent(userID, fileID, versionNumber)
	if err != nil {
		return err
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	_, copyErr := io.Copy(w, reader)
	return copyErr
}

// ListFileVersions handles GET /{id}/versions.
// ListFileVersions godoc
// @Summary      Lists file versions
// @Description  Retrieves a list of all versions of a file.
// @Tags         files
// @Produce      json
// @Param        id    path      string  true  "Target file ID"
// @Success      200  {array}  int  "List of version numbers"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /files/{id}/versions [get]
func (c *FileController) ListFileVersions(w http.ResponseWriter, r *http.Request) error {
	userID := middleware.MustUserID(r)

	fileID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	versions, err := c.fileService.ListVersions(userID, fileID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(versions)
}
