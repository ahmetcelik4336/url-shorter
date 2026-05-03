package services

import (
	"2/ent"
	"2/repositories"
	"errors"
	dto "shared/models"
)

type LogService interface {
	Create(urlId int, request dto.LogRequest) (*ent.Logs, error)
}

type logService struct {
	repo repositories.LogRepository
}

func NewLogService(r repositories.LogRepository) LogService {
	return &logService{repo: r}
}

func (s *logService) Create(urlId int, request dto.LogRequest) (*ent.Logs, error) {

	log, err := s.repo.Create(urlId, request)
	if err != nil {
		return nil, errors.New("hata oluştu")
	}

	return log, nil
}
