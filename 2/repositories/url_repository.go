package repositories

import (
	"context"
	dto "shared/models"
	"strings"

	"2/ent"
	"2/ent/url"
	"2/ent/user"
	"2/utils/jwtutil"

	"entgo.io/ent/dialect/sql"
)

type UrlRepository interface {
	FindByShortCode(short_code string) (*ent.Url, error)
	Create(userID int, short_code string, input dto.CreateUrlRequest) (*ent.Url, error)
	History(userID int) ([]*ent.Url, error)
	Update(input dto.UpdateUrlRequest) (*ent.Url, error)
	FindShortCodeOrAlias(shortcodeOrAlias string) (*ent.Url, error)
	FindByPassword(id int, password string) (*ent.Url, error)
	CreateBulk(userID int, inputs []dto.CreateUrlRequest) ([]*ent.Url, error)
}

type urlRepository struct {
	db *ent.Client
}

func NewUrlRepository(db *ent.Client) UrlRepository {
	return &urlRepository{
		db: db,
	}
}

func (r *urlRepository) FindByShortCode(short_code string) (*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.ShortCodeEQ(short_code)).
		//WithUser().
		Only(context.Background())
}

func (r *urlRepository) Create(userID int, short_code string, input dto.CreateUrlRequest) (*ent.Url, error) {

	query := r.db.Url.Create().
		SetLongURL(input.LongUrl).
		SetShortCode(short_code).
		SetUserID(userID).
		SetAlias(input.Alias).
		SetIsEncrypted(input.Password != "")
	if !input.ExpirationDate.IsZero() {
		query.SetExpirationDate(input.ExpirationDate)
	}

	return query.Save(context.Background())
}

func (r *urlRepository) CreateBulk(userID int, inputs []dto.CreateUrlRequest) ([]*ent.Url, error) {
	// 1. MapCreateBulk'ı slice ve bir callback fonksiyonu ile çağırıyoruz
	return r.db.Url.MapCreateBulk(inputs, func(b *ent.URLCreate, i int) {
		// 'i' o anki slice index'ini verir
		input := inputs[i]

		// Her kayıt için benzersiz kod üretimi (Önemli!)
		shortCode := jwtutil.GenerateShortID(5)

		// Builder'ı dolduruyoruz
		b.SetLongURL(input.LongUrl).
			SetShortCode(shortCode).
			SetUserID(userID).
			SetAlias(input.Alias).
			SetIsEncrypted(strings.TrimSpace(input.Password) != "").
			SetPassword(input.Password)

		if !input.ExpirationDate.IsZero() {
			b.SetExpirationDate(input.ExpirationDate)
		}
	}).Save(context.Background())
}

func (r *urlRepository) History(userID int) ([]*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.HasUserWith(user.IDEQ(userID))).
		Order(url.ByCreatedAt(sql.OrderDesc())).
		All(context.Background())

}

func (r *urlRepository) Update(input dto.UpdateUrlRequest) (*ent.Url, error) {

	return r.db.Url.
		UpdateOneID(input.ID).     // Güncellenecek kaydın ID'si
		SetLongURL(input.LongUrl). // Sadece değişmesini istediğin alanları setle
		Save(context.Background())
}

func (r *urlRepository) FindShortCodeOrAlias(shortcodeOrAlias string) (*ent.Url, error) {
	return r.db.Url.
		Query().
		Where(
			url.Or(
				url.ShortCodeEQ(shortcodeOrAlias),
				url.AliasEQ(shortcodeOrAlias),
			),
		).
		Only(context.Background())
}
func (r *urlRepository) FindByPassword(id int, password string) (*ent.Url, error) {
	return r.db.Url.
		Query().
		Where(url.PasswordEQ(password)).
		Where(url.IDEQ(id)).
		Only(context.Background())
}
