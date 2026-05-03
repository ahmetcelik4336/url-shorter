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

func NewContainer(db *ent.Client) *Container {

	userRepo := repositories.NewUserRepository(db)
	urlRepo := repositories.NewUrlRepository(db)
	analysisRepo := repositories.NewAnalysisRepository(db)
	logRepo := repositories.NewLogRepository(db)

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
