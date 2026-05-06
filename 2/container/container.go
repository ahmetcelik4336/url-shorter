package container

import (
	"2/ent"
	"2/repositories"
	"2/services"
)

type Container struct {
	UserService     services.UserService
	UrlService      services.UrlService
	AnalysisService services.AnalysisService
	LogService      services.LogService
}

func NewContainer(db *ent.Client, driver string) *Container {

	userRepo := repositories.NewUserRepository(db, driver)
	urlRepo := repositories.NewUrlRepository(db, driver)
	analysisRepo := repositories.NewAnalysisRepository(db, driver)
	logRepo := repositories.NewLogRepository(db, driver)

	userService := services.NewUserService(userRepo)
	urlService := services.NewUrlService(urlRepo)
	analysisService := services.NewAnalysisService(analysisRepo)
	logService := services.NewLogService(logRepo)

	return &Container{
		UserService:     userService,
		UrlService:      urlService,
		AnalysisService: analysisService,
		LogService:      logService,
	}
}
