package repositories

import (
	"context"
	dto "shared/models"
	"shared/utils"
	"time"

	"api/ent"
	"api/ent/url"
	"api/ent/user"
)

type UrlRepository interface {
	FindByShortCode(short_code string) (*ent.Url, error)
	Create(userID int, short_code string, input dto.CreateUrlRequest) (*ent.Url, error)
	History(userID int) ([]*ent.Url, error)
	HistoryById(userID int, id int) (*ent.Url, error)
	Update(input dto.UpdateUrlRequest) (*ent.Url, error)
	FindShortCodeOrAlias(shortcodeOrAlias string) (*ent.Url, error)
	FindByPassword(id int, password string) (*ent.Url, error)
	CreateBulk(userID int, inputs []dto.CreateUrlRequest) ([]*ent.Url, error)
	FindByAlias(alias string) (*ent.Url, error)
	FindByAliasWithId(id int, alias string) (*ent.Url, error)
	Delete(id int) error
}

type urlRepository struct {
	db      *ent.Client
	dialect string
}

func NewUrlRepository(db *ent.Client, dialect string) UrlRepository {
	return &urlRepository{
		db:      db,
		dialect: dialect,
	}
}
func (r *urlRepository) Delete(id int) error {
	return r.db.Url.DeleteOneID(id).Exec(context.Background())
}
func (r *urlRepository) FindByShortCode(short_code string) (*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.ShortCodeEQ(short_code)).
		//WithUser().
		Only(context.Background())
}

func (r *urlRepository) FindByAlias(alias string) (*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.AliasEQ(alias)).
		First(context.Background())
}

func (r *urlRepository) FindByAliasWithId(id int, alias string) (*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.AliasEQ(alias)).
		Where(url.IDNEQ(id)).
		First(context.Background())
}

func (r *urlRepository) Create(userID int, short_code string, input dto.CreateUrlRequest) (*ent.Url, error) {

	query := r.db.Url.Create().
		SetCreatedAt(time.Now()).
		SetLongURL(input.LongUrl).
		SetShortCode(short_code).
		SetUserID(userID).
		SetIsEncrypted(input.Password != "").
		SetPassword(input.Password)

	// ALIAS KONTROLÜ (Çökmeyi düzelten kısım):
	// Eğer alias girilmişse, adresini (&) alarak Ent'in beklediği *string tipinde gönderiyoruz.
	if input.Alias != "" {
		query.SetAlias(input.Alias)
	}

	// EXPIRATION DATE KONTROLÜ:
	// Şemanda bu da Nillable() olduğu için bunu da pointer olarak geçmeliyiz.
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
		shortCode := utils.GenerateShortID(5)

		// Builder'ı dolduruyoruz
		b.SetLongURL(input.LongUrl).
			SetShortCode(shortCode).
			SetUserID(userID).
			SetIsEncrypted(input.Password != "").
			SetPassword(input.Password)

		if input.Alias != "" {
			b.SetAlias(input.Alias)
		}

		if !input.ExpirationDate.IsZero() {
			b.SetExpirationDate(input.ExpirationDate)
		}

	}).Save(context.Background())
}

func (r *urlRepository) History(userID int) ([]*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.HasUserWith(user.IDEQ(userID))).
		Order(ent.Desc(url.FieldCreatedAt)).
		All(context.Background())

}

func (r *urlRepository) HistoryById(userID int, id int) (*ent.Url, error) {

	return r.db.Url.
		Query().
		Where(url.HasUserWith(user.IDEQ(userID))).
		Where(url.IDEQ(id)).
		First(context.Background())

}

func (r *urlRepository) Update(input dto.UpdateUrlRequest) (*ent.Url, error) {

	query := r.db.Url.
		UpdateOneID(input.ID).     // Güncellenecek kaydın ID'si
		SetLongURL(input.LongUrl). // Sadece değişmesini istediğin alanları setle
		SetIsEncrypted(input.Password != "").
		SetPassword(input.Password)

	if input.Alias != "" {
		query.SetAlias(input.Alias)
	}

	if !input.ExpirationDate.IsZero() {
		query.SetExpirationDate(input.ExpirationDate)
	}

	return query.Save(context.Background())
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
		First(context.Background())
}
func (r *urlRepository) FindByPassword(id int, password string) (*ent.Url, error) {
	return r.db.Url.
		Query().
		Where(url.PasswordEQ(password)).
		Where(url.IDEQ(id)).
		Only(context.Background())
}
