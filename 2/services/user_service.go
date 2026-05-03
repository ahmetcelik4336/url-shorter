package services

import (
	"2/repositories"
	"2/utils/jwtutil"
	"errors"
	dto "shared/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(email string, password string) (*dto.LoginResponse, error)
	Register(input dto.RegisterRequest) (*dto.RegisterResponse, error)
	EmailExists(email string) bool
	FindUser(id int) (*dto.UserResponse, error)
	ValidateToken(token string) (*dto.UserResponse, error)
	GetUsersCount() (*dto.UrlCountAnalysisResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) Login(email, password string) (*dto.LoginResponse, error) {

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. password compare (hash check)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("token error")
	}

	return dto.ToLoginResponse(token), nil
}

func (s *userService) EmailExists(email string) bool {

	_, err := s.repo.FindByEmail(email)
	if err != nil {
		return false
	}

	return true
}

func (s *userService) FindUser(id int) (*dto.UserResponse, error) {

	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}
func (s *userService) Register(input dto.RegisterRequest) (*dto.RegisterResponse, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Register(hash, input)
	if err != nil {
		return nil, err
	}
	var message string

	if user != nil {
		message = "Kayıt başarılı"
	} else {
		message = ""
	}
	return dto.ToRegisterResponse(user != nil, message), nil
}

func (s *userService) ValidateToken(token string) (*dto.UserResponse, error) {
	claims, err := jwtutil.ValidateToken(token)
	if err != nil {
		return nil, errors.New("Hata oluştu!")
	}
	user, err2 := s.FindUser(claims.UserID)
	if err2 != nil {
		return nil, errors.New("Kullanıcı bulunamadı!")
	}
	return user, nil
}

func (s *userService) GetUsersCount() (*dto.UrlCountAnalysisResponse, error) {

	user, err2 := s.repo.GetUsersCount()

	if err2 != nil {
		return nil, errors.New("Kullanıcı bulunamadı!")
	}

	return user, nil
}
