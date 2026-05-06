package repositories

import (
	"context"
	dto "shared/models"

	"2/ent"

	"entgo.io/ent/dialect/sql"
)

type LogRepository interface {
	Create(urlId int, request dto.LogRequest) (*ent.Logs, error)
	GetPerDayClick() (float64, error)
}

type logRepository struct {
	db *ent.Client
}

func NewLogRepository(db *ent.Client) LogRepository {
	return &logRepository{
		db: db,
	}
}
func (r *logRepository) GetPerDayClick() (float64, error) {
	var v []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	err := r.db.Logs.Query().
		Modify(func(s *sql.Selector) {
			s.Select(
				sql.As("DATE_FORMAT(created_at, '%Y-%m-%d')", "date"),
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
