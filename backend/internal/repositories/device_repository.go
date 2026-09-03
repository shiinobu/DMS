package repositories

import (
	"context"
	"errors"
	"time"

	"dms/backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDeviceNotFound = errors.New("device not found")

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) Create(
	ctx context.Context,
	device *models.Device,
) error {
	query := `
		INSERT INTO devices (
			device_id,
			device_name,
			serial_number,
			os_version,
			ip_address,
			location,
			status
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		device.DeviceID,
		device.DeviceName,
		device.SerialNumber,
		device.OSVersion,
		device.IPAddress,
		device.Location,
		device.Status,
	).Scan(
		&device.ID,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	return err
}

func (r *DeviceRepository) FindAll(
	ctx context.Context,
) ([]models.Device, error) {
	query := `
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
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
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

func (r *DeviceRepository) FindByDeviceID(
	ctx context.Context,
	deviceID string,
) (*models.Device, error) {
	query := `
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
		WHERE device_id = $1
	`

	var device models.Device

	err := r.db.QueryRow(
		ctx,
		query,
		deviceID,
	).Scan(
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

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrDeviceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *DeviceRepository) Update(
	ctx context.Context,
	deviceID string,
	device *models.Device,
) error {
	query := `
		UPDATE devices
		SET
			device_name = $1,
			serial_number = $2,
			os_version = $3,
			ip_address = $4,
			location = $5,
			updated_at = NOW()
		WHERE device_id = $6
		RETURNING
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
	`

	err := r.db.QueryRow(
		ctx,
		query,
		device.DeviceName,
		device.SerialNumber,
		device.OSVersion,
		device.IPAddress,
		device.Location,
		deviceID,
	).Scan(
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

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrDeviceNotFound
	}

	return err
}

func (r *DeviceRepository) Delete(
	ctx context.Context,
	deviceID string,
) error {
	query := `
		DELETE FROM devices
		WHERE device_id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		deviceID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

func (r *DeviceRepository) Heartbeat(
	ctx context.Context,
	deviceID string,
	ipAddress *string,
) (*models.Device, models.DeviceStatus, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, "", err
	}

	defer tx.Rollback(ctx)

	var previousStatus models.DeviceStatus

	err = tx.QueryRow(
		ctx,
		`
		SELECT status
		FROM devices
		WHERE device_id = $1
		FOR UPDATE
		`,
		deviceID,
	).Scan(&previousStatus)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", ErrDeviceNotFound
	}

	if err != nil {
		return nil, "", err
	}

	query := `
		UPDATE devices
		SET
			last_seen = NOW(),
			status = 'ONLINE',
			ip_address = COALESCE($2, ip_address),
			last_online_at = CASE
				WHEN status = 'OFFLINE' THEN NOW()
				ELSE last_online_at
			END,
			updated_at = NOW()
		WHERE device_id = $1
		RETURNING
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
	`

	var device models.Device

	err = tx.QueryRow(
		ctx,
		query,
		deviceID,
		ipAddress,
	).Scan(
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
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return &device, previousStatus, nil
}

func (r *DeviceRepository) FindOfflineCandidates(
	ctx context.Context,
	timeout time.Duration,
) ([]models.Device, error) {

	query := `
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
		WHERE status = 'ONLINE'
		  AND last_seen IS NOT NULL
		  AND last_seen < NOW() - $1::interval
	`

	rows, err := r.db.Query(
		ctx,
		query,
		timeout.String(),
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

func (r *DeviceRepository) MarkOffline(
	ctx context.Context,
	deviceID string,
) (*models.Device, error) {

	query := `
		UPDATE devices
		SET
			status = 'OFFLINE',
			last_offline_at = NOW(),
			updated_at = NOW()
		WHERE device_id = $1
		AND status = 'ONLINE'
		RETURNING
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
	`

	var device models.Device

	err := r.db.QueryRow(
		ctx,
		query,
		deviceID,
	).Scan(
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

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}
