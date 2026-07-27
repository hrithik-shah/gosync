package payload

import (
	"gosync/internal/api/dto"
	"time"
)

type DeviceInfo struct {
	ID         string    `json:"id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	LastSyncAt time.Time `json:"last_sync_at" validate:"required"`
}

type CreateDeviceRequest struct {
	Name string `json:"name" validate:"required"`
}

type ListDevicesResponse struct {
	Devices []DeviceInfo `json:"devices" validate:"required,dive"`
}

type GetDeviceResponse = DeviceInfo

func FromDeviceDTO(deviceDTO dto.DeviceDTO) DeviceInfo {
	return DeviceInfo{
		ID:         deviceDTO.ID,
		Name:       deviceDTO.Name,
		CreatedAt:  deviceDTO.CreatedAt,
		LastSyncAt: deviceDTO.LastSyncAt,
	}
}

func FromDeviceDTOSlice(deviceDTOs []dto.DeviceDTO) []DeviceInfo {
	devices := make([]DeviceInfo, len(deviceDTOs))
	for i, deviceDTO := range deviceDTOs {
		devices[i] = FromDeviceDTO(deviceDTO)
	}
	return devices
}
