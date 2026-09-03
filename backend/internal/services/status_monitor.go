package services

import (
	"context"
	"log"
	"time"

	"dms/backend/internal/repositories"
)

type StatusMonitor struct {
	repository      *repositories.DeviceRepository
	timeout         time.Duration
	checkInterval   time.Duration
	statusChanged   chan DeviceStatusChangedEvent
}

func NewStatusMonitor(
	repository *repositories.DeviceRepository,
	timeout time.Duration,
	checkInterval time.Duration,
) *StatusMonitor {

	return &StatusMonitor{
		repository:    repository,
		timeout:       timeout,
		checkInterval: checkInterval,
		statusChanged: make(chan DeviceStatusChangedEvent, 100),
	}
}

func (m *StatusMonitor) Start(ctx context.Context) {

	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()

	log.Println("Device status monitor started")

	for {
		select {

		case <-ticker.C:
			if err := m.checkDevices(ctx); err != nil {
				log.Println("status monitor error:", err)
			}

		case <-ctx.Done():
			log.Println("Device status monitor stopped")
			return
		}
	}
}

func (m *StatusMonitor) checkDevices(
	ctx context.Context,
) error {

	devices, err := m.repository.FindOfflineCandidates(
		ctx,
		m.timeout,
	)

	if err != nil {
		return err
	}

	for _, candidate := range devices {

		device, err := m.repository.MarkOffline(
			ctx,
			candidate.DeviceID,
		)

		if err != nil {

			log.Printf(
				"failed to mark device %s offline: %v",
				candidate.DeviceID,
				err,
			)

			continue
		}

		// Device may have received heartbeat
		// between FindOfflineCandidates and MarkOffline.
		if device == nil {
			continue
		}

		log.Printf(
			"Device %s changed status: ONLINE -> OFFLINE",
			device.DeviceID,
		)

		m.statusChanged <- DeviceStatusChangedEvent{
			Type:     "DEVICE_STATUS_CHANGED",
			DeviceID: device.DeviceID,
			Status:   device.Status,
			Device:   *device,
		}
	}

	return nil
}

func (m *StatusMonitor) StatusChanged() <-chan DeviceStatusChangedEvent {

	return m.statusChanged
}