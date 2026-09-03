package services

import (
	"context"
	"log"
	"time"

	"dms/backend/internal/repositories"
)

type StatusMonitor struct {
	repository    *repositories.DeviceRepository
	timeout       time.Duration
	checkInterval time.Duration
	statusChanged chan DeviceStatusChangedEvent
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

	log.Println("Device Status Monitor Started")

	for {
		select {
		case <-ticker.C:
			if err := m.checkDevices(ctx); err != nil {
				log.Printf(
					"Status Monitor Error: %v",
					err,
				)
			}

		case <-ctx.Done():
			log.Println("Device Status Monitor Stopped")
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
				"Failed to Mark Device %s Offline: %v",
				candidate.DeviceID,
				err,
			)

			continue
		}

		if device == nil {
			continue
		}

		log.Printf(
			"Device %s Changed Status: ONLINE -> OFFLINE",
			device.DeviceID,
		)

		m.statusChanged <- DeviceStatusChangedEvent{
			Type:     DeviceStatusChangedEventType,
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
