package services

import (
	"api/ent"
	"api/internal/repositories"
	"encoding/json"
	"errors"
	dto "shared/models"
)

type SettingService interface {
	Save(key string, data any) (*ent.Setting, error)
	GetGeneralSettings() (*dto.GeneralSettings, error)
}

type settingService struct {
	repo repositories.SettingRepository
}

func NewSettingService(r repositories.SettingRepository) SettingService {
	return &settingService{repo: r}
}
func (s *settingService) Save(key string, data any) (*ent.Setting, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("json hata oluştu")
	}

	setting, err := s.repo.Save(dto.SettingResponse{
		Key:     key,
		Content: string(jsonData),
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return setting, nil
}
func (s *settingService) GetGeneralSettings() (*dto.GeneralSettings, error) {

	settings, err := s.repo.GetGeneralSettings()
	if err != nil {
		return nil, errors.New("hata oluştu")
	}
	return settings, nil
}
