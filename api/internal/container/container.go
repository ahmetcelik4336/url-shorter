package container

import (
	"api/ent"
	"api/internal/repositories"
	"api/internal/services"
)

type Container struct {
	UserService     services.UserService
	UrlService      services.UrlService
	AnalysisService services.AnalysisService
	LogService      services.LogService
	SettingService  services.SettingService
}

func NewContainer(db *ent.Client, driver string) *Container {

	userRepo := repositories.NewUserRepository(db, driver)
	urlRepo := repositories.NewUrlRepository(db, driver)
	analysisRepo := repositories.NewAnalysisRepository(db, driver)
	logRepo := repositories.NewLogRepository(db, driver)
	settingsRepo := repositories.NewSettingRepository(db, driver)

	userService := services.NewUserService(userRepo)
	urlService := services.NewUrlService(urlRepo)
	analysisService := services.NewAnalysisService(analysisRepo)
	logService := services.NewLogService(logRepo)
	settingService := services.NewSettingService(settingsRepo)

	return &Container{
		UserService:     userService,
		UrlService:      urlService,
		AnalysisService: analysisService,
		LogService:      logService,
		SettingService:  settingService,
	}
}
