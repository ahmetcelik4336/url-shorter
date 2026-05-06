package repositories

import (
	"context"
	dto "shared/models"

	"2/ent"
	"2/ent/user"
)

type UserRepository interface {
	FindByEmail(email string) (*ent.User, error)
	Register(hash []byte, input dto.RegisterRequest) (*ent.User, error)
	FindById(id int) (*ent.User, error)
	GetUsersCount() (*dto.UrlCountAnalysisResponse, error)
}

type userRepository struct {
	db      *ent.Client
	dialect string
}

func NewUserRepository(db *ent.Client, dialect string) UserRepository {
	return &userRepository{
		db:      db,
		dialect: dialect,
	}
}

func (r *userRepository) FindByEmail(email string) (*ent.User, error) {

	return r.db.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(context.Background())
}

func (r *userRepository) FindById(id int) (*ent.User, error) {

	return r.db.User.
		Query().
		Where(user.IDEQ(id)).
		Only(context.Background())
}

func (r *userRepository) Register(hash []byte, input dto.RegisterRequest) (*ent.User, error) {

	return r.db.User.
		Create().
		SetAd(input.Name).
		SetEmail(input.Email).
		SetPassword(string(hash)).
		Save(context.Background())
}

func (r *userRepository) GetUsersCount() (*dto.UrlCountAnalysisResponse, error) {
	q := r.db.User.
		Query()

	count, err := q.Count(context.Background())

	return &dto.UrlCountAnalysisResponse{
		Count: count,
	}, err
}
