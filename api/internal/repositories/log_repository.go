package repositories

import (
	"context"
	"fmt"
	"shared/db"
	dto "shared/models"
	"time"

	"api/ent"
	"api/ent/logs"
	"api/ent/url"
	"api/ent/user"

	"entgo.io/ent/dialect/sql"
)

type LogRepository interface {
	Create(urlId, userId int, request dto.LogRequest) (*ent.Logs, error)
	GetPerDayClick() (float64, error)
	TotalReading(userID int) (*dto.UrlTrackAnalysisResponse, error)
	TotalReadingQr(userID int) (*dto.UrlTrackAnalysisResponse, error)
	TotalReadingUrl(userID int) (*dto.UrlTrackAnalysisResponse, error)
	GetLastReading(userID int) (*dto.GetLastReadingResponse, error)
	GetUrlTrackAnalysis(userID int, request *dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error)
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
			groupByExpr := dialectHandler.GroupByDate(logs.FieldCreatedAt)
			s.Select(
				sql.As(dateExpr, "date"),
				sql.As("COUNT(*)", "count"),
			).
				GroupBy(groupByExpr).
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
func (r *logRepository) Create(urlId, userId int, request dto.LogRequest) (*ent.Logs, error) {
	q := r.db.Logs.
		Create().
		SetDevice(request.Device).
		SetIP(request.Ip).
		SetReferer(request.Referer).
		SetLogID(urlId).
		SetType(request.Type).
		SetCreatedAt(time.Now())

	if userId > 0 {
		q.SetUserID(userId)
	}

	return q.Save(context.Background())
}

func (r *logRepository) TotalReading(userID int) (*dto.UrlTrackAnalysisResponse, error) {
	totalReading, _ := r.db.Logs.
		Query().
		Where(logs.HasUserWith(user.IDEQ(userID))).
		Count(context.Background())

	return &dto.UrlTrackAnalysisResponse{
		Count: totalReading,
		Type:  "totalReading",
		Title: "Total reading",
	}, nil
}

func (r *logRepository) TotalReadingQr(userID int) (*dto.UrlTrackAnalysisResponse, error) {
	totalReading, _ := r.db.Logs.
		Query().
		Where(logs.HasUserWith(user.IDEQ(userID))).
		Where(logs.TypeEQ("qr")).
		Count(context.Background())

	return &dto.UrlTrackAnalysisResponse{
		Count: totalReading,
		Type:  "totalQrReading",
		Title: "Total Qr reading",
	}, nil
}
func (r *logRepository) TotalReadingUrl(userID int) (*dto.UrlTrackAnalysisResponse, error) {
	totalReading, _ := r.db.Logs.
		Query().
		Where(logs.HasUserWith(user.IDEQ(userID))).
		Where(logs.TypeEQ("url")).
		Count(context.Background())

	return &dto.UrlTrackAnalysisResponse{
		Count: totalReading,
		Type:  "totalUrlReading",
		Title: "Total Url Reading",
	}, nil
}
func (r *logRepository) GetLastReading(userID int) (*dto.GetLastReadingResponse, error) {
	lastLog, err := r.db.Logs.
		Query().
		Where(logs.HasUserWith(user.IDEQ(userID))).
		Order(ent.Desc(logs.FieldCreatedAt)).
		First(context.Background())

	if err != nil {
		return nil, err
	}

	// Geçen süreyi hesapla
	diff := time.Since(lastLog.CreatedAt)
	var timeAgo string

	// Basit bir Türkçe zaman farkı hesaplayıcı
	switch {
	case diff < time.Minute:
		timeAgo = "Biraz önce"
	case diff < time.Hour:
		timeAgo = fmt.Sprintf("%.0f dakika önce", diff.Minutes())
	case diff < 24*time.Hour:
		timeAgo = fmt.Sprintf("%.0f saat önce", diff.Hours())
	default:
		timeAgo = fmt.Sprintf("%.0f gün önce", diff.Hours()/24)
	}

	return &dto.GetLastReadingResponse{
		LastReading: timeAgo, // Örn: "2 saat önce"
	}, nil
}

func (r *logRepository) GetUrlTrackAnalysis(userID int, request *dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error) {
	var v []*dto.UsageAnalysisResponse

	q := r.db.Logs.Query().
		Where(logs.HasUserWith(user.IDEQ(userID)))

	if !request.Start.IsZero() {
		q = q.Where(
			logs.CreatedAtGTE(request.Start),
		)
	} else {
		q = q.Where(
			logs.CreatedAtGTE(time.Now().Add(-7 * 24 * time.Hour)),
		)
	}

	if !request.End.IsZero() {
		q = q.Where(
			logs.CreatedAtLTE(request.End),
		)
	} else {
		q = q.Where(
			logs.CreatedAtLTE(time.Now()),
		)
	}

	err := q.Modify(func(s *sql.Selector) {

		dialectHandler := db.GetDialect(r.dailect)
		dateExpr := dialectHandler.FormatDate(url.FieldCreatedAt)

		s.Select(
			sql.As(dateExpr, "date"),
			sql.As("COUNT(*)", "count"),
		).
			GroupBy(dateExpr).
			OrderBy(sql.Desc("date"))
	}).Scan(context.Background(), &v)

	return v, err
}
