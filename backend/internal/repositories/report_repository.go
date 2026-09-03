package repositories

import (
	"context"
	"errors"

	"dms/backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository interface {
	GetSummary(ctx context.Context) (*models.ReportSummary, error)
	GetDevicesForExport(ctx context.Context) ([]models.Device, error)
}

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(
	db *pgxpool.Pool,
) ReportRepository {
	return &reportRepository{
		db: db,
	}
}

func (r *reportRepository) GetSummary(
	ctx context.Context,
) (*models.ReportSummary, error) {
	var result models.ReportSummary

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'ONLINE'),
			COUNT(*) FILTER (WHERE status = 'OFFLINE')
		FROM devices
		`,
	).Scan(
		&result.TotalDevices,
		&result.OnlineDevices,
		&result.OfflineDevices,
	)

	if err != nil {
		return nil, err
	}

	var lastOnline models.LastOnline

	err = r.db.QueryRow(
		ctx,
		`
		SELECT
			device_id,
			device_name,
			last_online_at
		FROM devices
		WHERE last_online_at IS NOT NULL
		ORDER BY last_online_at DESC
		LIMIT 1
		`,
	).Scan(
		&lastOnline.DeviceID,
		&lastOnline.DeviceName,
		&lastOnline.LastOnlineAt,
	)

	if err == nil {
		result.LastOnline = &lastOnline
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var lastOffline models.LastOffline

	err = r.db.QueryRow(
		ctx,
		`
		SELECT
			device_id,
			device_name,
			last_offline_at
		FROM devices
		WHERE last_offline_at IS NOT NULL
		ORDER BY last_offline_at DESC
		LIMIT 1
		`,
	).Scan(
		&lastOffline.DeviceID,
		&lastOffline.DeviceName,
		&lastOffline.LastOfflineAt,
	)

	if err == nil {
		result.LastOffline = &lastOffline
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &result, nil
}

func (r *reportRepository) GetDevicesForExport(
    ctx context.Context,
) ([]models.Device, error) {
    rows, err := r.db.Query(
        ctx,
        `
        SELECT
            id,
            device_id,
            device_name,
            serial_number,
            os_version,
            ip_address::text,
            location,
            status,
            last_seen,
            last_online_at,
            last_offline_at,
            created_at,
            updated_at
        FROM devices
        ORDER BY id ASC
        `,
    )

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    devices := make([]models.Device, 0)

    for rows.Next() {
        var device models.Device

        err := rows.Scan(
            &device.ID,
            &device.DeviceID,
            &device.DeviceName,
            &device.SerialNumber,
            &device.OSVersion,
            &device.IPAddress,
            &device.Location,
            &device.Status,
            &device.LastSeen,
            &device.LastOnlineAt,
            &device.LastOfflineAt,
            &device.CreatedAt,
            &device.UpdatedAt,
        )

        if err != nil {
            return nil, err
        }

        devices = append(devices, device)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return devices, nil
}