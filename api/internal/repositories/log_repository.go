package repositories

import (
	"context"
	"shared/db"
	dto "shared/models"

	"api/ent"
	"api/ent/logs"

	"entgo.io/ent/dialect/sql"
)

type LogRepository interface {
	Create(urlId int, request dto.LogRequest) (*ent.Logs, error)
	GetPerDayClick() (float64, error)
}

type logRepository struct {
	db      *ent.Client
	dailect string
}

func NewLogRepository(db *ent.Client, dialect string) LogRepository {
	return &logRepository{
		db:      db,
		dailect: dialect,
	}
}
func (r *logRepository) GetPerDayClick() (float64, error) {
	var v []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	err := r.db.Logs.Query().
		Modify(func(s *sql.Selector) {

			dialectHandler := db.GetDialect(r.dailect)
			dateExpr := dialectHandler.FormatDate(logs.FieldCreatedAt)

			s.Select(
				sql.As(dateExpr, "date"),
				sql.As("COUNT(*)", "count"),
			).
				GroupBy("DATE(created_at)").
				OrderBy(sql.Desc("date"))
		}).
		Scan(context.Background(), &v)

	if err != nil {
		return 0, err
	}

	if len(v) == 0 {
		return 0, nil
	}

	var totalCount int
	for _, item := range v {
		totalCount += item.Count
	}

	average := float64(totalCount) / float64(len(v))
	return average, nil
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
