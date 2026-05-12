package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	dto "shared/models"

	"api/ent"
	"api/ent/setting"
)

type SettingRepository interface {
	Save(request dto.SettingResponse) (*ent.Setting, error)
	GetGeneralSettings() (*dto.GeneralSettings, error)
}

type settingRepository struct {
	db      *ent.Client
	dialect string
}

func NewSettingRepository(db *ent.Client, dialect string) SettingRepository {
	return &settingRepository{
		db:      db,
		dialect: dialect,
	}
}

func (r *settingRepository) Save(request dto.SettingResponse) (*ent.Setting, error) {
	return r.db.Setting.
		Create().
		SetSettingsKey(request.Key).
		SetSettingContent(request.Content).
		Save(context.Background())
}

func (r *settingRepository) GetGeneralSettings() (*dto.GeneralSettings, error) {

	data, err := r.db.Setting.
		Query().
		Where(setting.SettingsKeyEQ("general")).
		First(context.Background())

	if err != nil {
		return nil, fmt.Errorf("ayar bulunamadı: %v", err)
	}

	var result dto.GeneralSettings

	err = json.Unmarshal([]byte(data.SettingContent), &result)
	if err != nil {
		return nil, fmt.Errorf("json çözümleme hatası: %v", err)
	}

	return &result, nil

}
