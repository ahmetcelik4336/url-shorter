package repositories

import (
	"context"
	dto "shared/models"

	"2/ent"
)

type LogRepository interface {
	Create(urlId int, request dto.LogRequest) (*ent.Logs, error)
}

type logRepository struct {
	db *ent.Client
}

func NewLogRepository(db *ent.Client) LogRepository {
	return &logRepository{
		db: db,
	}
}

func (r *logRepository) Create(urlId int, request dto.LogRequest) (*ent.Logs, error) {
	return r.db.Logs.
		Create().
		SetDevice(request.Device).
		SetIP(request.Ip).
		SetReferer(request.Referer).
		SetLogID(urlId).
		Save(context.Background())
}
