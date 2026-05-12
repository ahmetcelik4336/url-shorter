package services

import (
	"api/internal/repositories"
	"errors"
	dto "shared/models"
)

type AnalysisService interface {
	GetURLStats(userID int, request dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error)
	GetUserURLCount(userID int, request dto.UrlCountAnalysisRequest) (*dto.UrlCountAnalysisResponse, error)
	GetURLCountTotal() (*dto.UrlCountAnalysisResponse, error)
}

type analysisService struct {
	repo repositories.AnalysisRepository
}

func NewAnalysisService(r repositories.AnalysisRepository) AnalysisService {
	return &analysisService{repo: r}
}
func (s *analysisService) GetURLCountTotal() (*dto.UrlCountAnalysisResponse, error) {
	res, err := s.repo.GetURLCountTotal()
	if err != nil {
		return nil, errors.New("hata oluştu")
	}

	return res, nil
}
func (s *analysisService) GetURLStats(userID int, request dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error) {

	res, err := s.repo.GetURLStats(userID, request)
	if err != nil {
		return nil, errors.New("hata oluştu")
	}

	if res == nil {
		return nil, errors.New("veri yok")
	}

	return res, nil
}

func (s *analysisService) GetUserURLCount(userID int, request dto.UrlCountAnalysisRequest) (*dto.UrlCountAnalysisResponse, error) {

	analysis, err := s.repo.GetUserURLCount(userID, request)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return analysis, nil
}
