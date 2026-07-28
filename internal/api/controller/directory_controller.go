package controller

import (
	"encoding/json"
	"net/http"

	"gosync/internal/api/apperror"
	"gosync/internal/api/middleware"
	"gosync/internal/api/payload"
	"gosync/internal/api/service"
	"gosync/internal/api/utils/httputil"
)

type DirectoryController struct {
	directoryService *service.DirectoryService
}

func NewDirectoryController() *DirectoryController {
	return &DirectoryController{directoryService: service.NewDirectoryService()}
}

// Create godoc
// @Description  Creates a new directory
// @Summary      Create a new directory
// @Tags         directories
// @Accept       json
// @Param        body  	body  payload.CreateDirectoryRequest  true  "Directory details"
// @Success      201  {object} 	payload.CreateDirectoryResponse
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories [post]
func (c *DirectoryController) Create(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	var req payload.CreateDirectoryRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	dirID, err := c.directoryService.Create(userID, req.Name, req.ParentDirectoryID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(payload.CreateDirectoryResponse{DirectoryID: dirID.String()})
}

// GetRootMetadata godoc
// @Description  Gets metadata for the root directory
// @Summary      Get root directory metadata
// @Tags         directories
// @Success      200  {object} 	payload.GetMetadataResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories/root [get]
func (c *DirectoryController) GetRootMetadata(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	directoryDTO, err := c.directoryService.GetRootMetadata(userID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(payload.GetMetadataResponse{
		ID:                directoryDTO.ID,
		Name:              directoryDTO.Name,
		ParentDirectoryID: directoryDTO.ParentDirectoryID,
		Hash:              directoryDTO.Hash,
	})
}

// GetMetadata godoc
// @Description  Gets metadata for a directory
// @Summary      Get directory metadata
// @Tags         directories
// @Param        id  	path  string  true  "Directory ID"
// @Success      200  {object} 	payload.GetMetadataResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories/{id} [get]
func (c *DirectoryController) GetMetadata(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	dirID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	directoryDTO, err := c.directoryService.GetMetadata(userID, dirID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(payload.GetMetadataResponse{
		ID:                directoryDTO.ID,
		Name:              directoryDTO.Name,
		ParentDirectoryID: directoryDTO.ParentDirectoryID,
		Hash:              directoryDTO.Hash,
	})
}

// ListContents godoc
// @Description  Lists contents of a directory
// @Summary      List directory contents
// @Tags         directories
// @Param        id  	path  string  true  "Directory ID"
// @Success      200  {object} 	payload.ListDirectoryContentsResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories/{id}/list [get]
func (c *DirectoryController) ListContents(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	dirID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	files, err := c.directoryService.ListContents(userID, dirID)
	if err != nil {
		return err
	}

	response := payload.ListDirectoryContentsResponse{
		Files: payload.FromFileDTOSlice(files),
	}

	return json.NewEncoder(w).Encode(response)
}

// Update godoc
// @Description  Updates a directory
// @Summary      Update a directory
// @Tags         directories
// @Accept       json
// @Param        id  	path  string  true  "Directory ID"
// @Param        body  	body  payload.UpdateDirectoryRequest  true  "New directory details"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories/{id} [patch]
func (c *DirectoryController) Update(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	dirID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	var req payload.UpdateDirectoryRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	err = c.directoryService.Update(userID, dirID, req.NewName)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Move godoc
// @Description  Moves a directory
// @Summary      Move a directory
// @Tags         directories
// @Accept       json
// @Param        id  	path  string  true  "Directory ID"
// @Param        body  	body  payload.MoveDirectoryRequest  true  "New directory details"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories/{id}/move [post]
func (c *DirectoryController) Move(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	dirID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	var req payload.MoveDirectoryRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	err = c.directoryService.Move(userID, dirID, req.NewParentDirectoryID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Delete godoc
// @Description  Deletes a directory
// @Summary      Delete a directory
// @Tags         directories
// @Param        id  	path  string  true  "Directory ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /directories/{id} [delete]
func (c *DirectoryController) Delete(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	dirID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	err = c.directoryService.Delete(userID, dirID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
