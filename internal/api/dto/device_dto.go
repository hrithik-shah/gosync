package dto

import "time"

type DeviceDTO struct {
	ID           string
	Name         string
	LastSyncAt   time.Time
	LastRootHash string
	CreatedAt    time.Time
}
