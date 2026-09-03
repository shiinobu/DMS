package services

import (
	"context"

	"dms/backend/internal/models"
	"dms/backend/internal/repositories"
)

type ReportService interface {
	GetSummary(ctx context.Context) (*models.ReportSummary, error)
	GetDevicesForExport(ctx context.Context) ([]models.Device, error)
}

type reportService struct {
	repository repositories.ReportRepository
}

func NewReportService(
	repository repositories.ReportRepository,
) ReportService {
	return &reportService{
		repository: repository,
	}
}

func (s *reportService) GetSummary(
	ctx context.Context,
) (*models.ReportSummary, error) {
	return s.repository.GetSummary(ctx)
}

func (s *reportService) GetDevicesForExport(
    ctx context.Context,
) ([]models.Device, error) {
    return s.repository.GetDevicesForExport(ctx)
}