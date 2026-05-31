package services

import (
	"api/ent"
	"api/internal/repositories"
	"errors"
	dto "shared/models"
	"time"
)

type LogService interface {
	Create(urlId, userId int, request dto.LogRequest) (*ent.Logs, error)
	GetPerDayClick() (*dto.ResponseClicked, error)
	TotalReading(userId int) (*dto.UrlTrackAnalysisResponseBatch, error)
	GetUrlTrackAnalysis(userId int, request *dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error)
	LogDatatable(draw, start, length, userId int, searchval, orderColumnName, orderDir string, startDate, endDate time.Time) (*dto.DataTableResponse, error)
}

type logService struct {
	repo repositories.LogRepository
}

func NewLogService(r repositories.LogRepository) LogService {
	return &logService{repo: r}
}
func (s *logService) GetPerDayClick() (*dto.ResponseClicked, error) {
	log, err := s.repo.GetPerDayClick()
	if err != nil {
		return nil, errors.New("hata oluştu")
	}

	return &dto.ResponseClicked{
		Status: true,
		Count:  int(log),
	}, nil
}
func (s *logService) Create(urlId, userId int, request dto.LogRequest) (*ent.Logs, error) {

	log, err := s.repo.Create(urlId, userId, request)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return log, nil
}

func (s *logService) TotalReading(userId int) (*dto.UrlTrackAnalysisResponseBatch, error) {
	totalReading, _ := s.repo.TotalReading(userId)
	totalReadingQr, _ := s.repo.TotalReadingQr(userId)
	totalReadingUrl, _ := s.repo.TotalReadingUrl(userId)
	getLastReading, _ := s.repo.GetLastReading(userId)

	return &dto.UrlTrackAnalysisResponseBatch{
		LastReading: getLastReading,
		Analysis: []*dto.UrlTrackAnalysisResponse{
			totalReading,
			totalReadingQr,
			totalReadingUrl,
		},
	}, nil
}

func (s *logService) GetUrlTrackAnalysis(userId int, request *dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error) {
	GetUrlTrackAnalysis, _ := s.repo.GetUrlTrackAnalysis(userId, request)
	return GetUrlTrackAnalysis, nil
}

func (s *logService) LogDatatable(draw, start, length, userId int, searchval, orderColumnName, orderDir string, startDate, endDate time.Time) (*dto.DataTableResponse, error) {

	data, totalRecords, err := s.repo.LogDatatable(start, length, userId, searchval, orderColumnName, orderDir, startDate, endDate)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &dto.DataTableResponse{
		Draw:            draw,
		RecordsTotal:    int64(totalRecords),
		RecordsFiltered: totalRecords,
		Data:            data,
	}, nil
}
