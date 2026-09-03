package services

import (
	"context"
	"strings"

	"dms/backend/internal/models"
	"dms/backend/internal/repositories"
)

type DeviceService struct {
	repository *repositories.DeviceRepository
	statusChanged chan DeviceStatusChangedEvent
}

func NewDeviceService(
	repository *repositories.DeviceRepository,
) *DeviceService {
	return &DeviceService{
		repository: repository,
		statusChanged: make(chan DeviceStatusChangedEvent, 100),
	}
}

func (s *DeviceService) Create(
	ctx context.Context,
	device *models.Device,
) error {
	device.DeviceID = strings.TrimSpace(device.DeviceID)
	device.DeviceName = strings.TrimSpace(device.DeviceName)
	device.SerialNumber = strings.TrimSpace(device.SerialNumber)

	device.Status = models.DeviceStatusOffline

	return s.repository.Create(ctx, device)
}

func (s *DeviceService) FindAll(
	ctx context.Context,
) ([]models.Device, error) {
	return s.repository.FindAll(ctx)
}

func (s *DeviceService) FindByDeviceID(
	ctx context.Context,
	deviceID string,
) (*models.Device, error) {
	return s.repository.FindByDeviceID(ctx, deviceID)
}

func (s *DeviceService) Update(
	ctx context.Context,
	deviceID string,
	device *models.Device,
) error {
	deviceID = strings.TrimSpace(deviceID)
	device.DeviceName = strings.TrimSpace(device.DeviceName)
	device.SerialNumber = strings.TrimSpace(device.SerialNumber)

	return s.repository.Update(
		ctx,
		deviceID,
		device,
	)
}

func (s *DeviceService) Delete(
	ctx context.Context,
	deviceID string,
) error {
	deviceID = strings.TrimSpace(deviceID)

	return s.repository.Delete(
		ctx,
		deviceID,
	)
}

func (s *DeviceService) Heartbeat(
	ctx context.Context,
	deviceID string,
	ipAddress *string,
) (*models.Device, models.DeviceStatus, error) {

	deviceID = strings.TrimSpace(deviceID)

	device, previousStatus, err := s.repository.Heartbeat(
		ctx,
		deviceID,
		ipAddress,
	)

	if err != nil {
		return nil, "", err
	}

	if previousStatus != device.Status {

		s.statusChanged <- DeviceStatusChangedEvent{
			Type:     "DEVICE_STATUS_CHANGED",
			DeviceID: device.DeviceID,
			Status:   device.Status,
			Device:   *device,
		}
	}

	return device, previousStatus, nil
}

func (s *DeviceService) StatusChanged() <-chan DeviceStatusChangedEvent {

	return s.statusChanged
}
