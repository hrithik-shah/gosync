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

type DeviceController struct {
	deviceService *service.DeviceService
}

func NewDeviceController() *DeviceController {
	return &DeviceController{deviceService: service.NewDeviceService()}
}

// Create godoc
// @Summary      Create a new device
// @Description  Creates a new device account
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        body  body      payload.CreateDeviceRequest  true  "Registration details"
// @Success      201   {object}  payload.RegisterResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /devices [post]
func (c *DeviceController) Create(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	var req payload.CreateDeviceRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	err := c.deviceService.Create(userID, req.Name)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// List godoc
// @Summary      List devices
// @Description  Retrieves a list of all devices
// @Tags         devices
// @Produce      json
// @Success      200   {object}  payload.ListDevicesResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /devices [get]
func (c *DeviceController) List(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	var req payload.LoginRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	devices, err := c.deviceService.List(userID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(payload.ListDevicesResponse{
		Devices: payload.FromDeviceDTOSlice(devices),
	})
}

// Get godoc
// @Summary      Get device
// @Description  Retrieves a device by ID
// @Tags         devices
// @Produce      json
// @Param        id  	path  string  true  "Device ID"
// @Success      200   {object}  payload.GetDeviceResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /devices/{id} [get]
func (c *DeviceController) Get(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	deviceID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		return err
	}

	device, err := c.deviceService.Get(userID, deviceID)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(payload.FromDeviceDTO(device))
}
