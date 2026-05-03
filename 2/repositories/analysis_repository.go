package repositories

import (
	"2/ent"
	"2/ent/url"
	"2/ent/user"
	"context"
	dto "shared/models"
	"shared/utils"

	"entgo.io/ent/dialect/sql"
)

type AnalysisRepository interface {
	GetUserURLCount(userID int, request dto.UrlCountAnalysisRequest) (*dto.UrlCountAnalysisResponse, error)
	GetURLStats(userID int, request dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error)
}

type analysisRepository struct {
	db *ent.Client
}

func NewAnalysisRepository(db *ent.Client) AnalysisRepository {
	return &analysisRepository{
		db: db,
	}
}

func (r *analysisRepository) GetUserURLCount(userID int, request dto.UrlCountAnalysisRequest) (*dto.UrlCountAnalysisResponse, error) {
	q := r.db.Url.
		Query().
		Where(url.HasUserWith(user.IDEQ(userID)))

	if request.Date != "" {
		start, end := utils.GenerateDates(request.Date)
		q.Where(url.CreatedAtGTE(start), url.CreatedAtLTE(end))
	}

	count, err := q.Count(context.Background())

	return &dto.UrlCountAnalysisResponse{
		Count: count,
	}, err
}

func (r *analysisRepository) GetURLStats(userID int, request dto.UsageAnalysisRequest) ([]*dto.UsageAnalysisResponse, error) {
	var v []*dto.UsageAnalysisResponse

	q := r.db.Url.Query().
		Where(url.HasUserWith(user.IDEQ(userID)))

	if request.Date != "" {
		start, end := utils.GenerateDates(request.Date)
		q = q.Where(
			url.CreatedAtGTE(start),
			url.CreatedAtLTE(end),
		)
	}

	err := q.
		Modify(func(s *sql.Selector) {
			s.Select(
				sql.As("DATE_FORMAT(created_at, '%Y-%m-%d')", "date"),
				sql.As("COUNT(*)", "count"),
			).
				GroupBy("DATE(created_at)").
				OrderBy(sql.Desc("date"))
		}).
		Scan(context.Background(), &v)

	return v, err
}

/*
func (r *analysisRepository) GetURLStats(userID int, request dto.UsageAnalysisRequest) ([]dto.UsageAnalysisResponse, error) {
	q := r.db.Url.Query().Where(url.HasUserWith(user.IDEQ(userID)))

	// Tarih filtresi varsa uygula
	if request.Date != "" {
		start, end := jwtutil.GenerateDates(request.Date)
		q = q.Where(
			url.CreatedAtGTE(start),
			url.CreatedAtLTE(end),
		)
	}

	urls, err := q.All(context.Background())
	if err != nil {
		return nil, err
	}

	statsMap := make(map[string]int)

	for _, u := range urls {
		// CreatedAt'i "YYYY-MM-DD" formatına çeviriyoruz
		dateStr := u.CreatedAt.Format("2006-01-02")
		statsMap[dateStr]++
	}

	// 3. Map'teki sonuçları istenen Response slice'ına dönüştürüyoruz
	var v []dto.UsageAnalysisResponse
	for date, total := range statsMap {
		v = append(v, dto.UsageAnalysisResponse{
			Date:  date,
			Count: total,
		})
	}

	return v, nil
}*/
