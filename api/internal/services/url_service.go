package services

import (
	"api/ent"
	"api/internal/repositories"
	"errors"
	netURL "net/url"
	dto "shared/models"
	"shared/utils"
	"strings"
	"time"
)

type UrlService interface {
	FindByShortCode(shortcode string) (*ent.Url, error)
	Create(userID int, req dto.CreateUrlRequest, type_ string) (*dto.GeneralResponse[dto.UrlResponse], error)
	History(userID int) ([]*dto.UrlResponse, error)
	Update(req dto.UpdateUrlRequest) (*dto.GeneralResponse[any], error)
	Redirect(shortcodeOrAlias, password string) (*dto.GeneralResponse[dto.UrlResponse], error)
	Bulkcreate(userID int, req []dto.CreateUrlRequest) (*dto.GeneralResponse[dto.UrlResponse], error)
	HistoryById(userID int, id int) (*dto.UrlResponse, error)
	Delete(id int) (*dto.GeneralResponse[any], error)
	GetById(id int) (*ent.Url, error)
}

type urlService struct {
	repo repositories.UrlRepository
}

func NewUrlService(r repositories.UrlRepository) UrlService {
	return &urlService{repo: r}
}

/*
func (s *urlService) FindByAlias(alias string) (*dto.GeneralResponse[any], error) {

		urlResultt, _ := s.repo.FindByAlias(alias)
		if urlResultt.ID > 0 {
			return &dto.GeneralResponse[any]{
				Status:  false,
				Message: "Takma ad kayıtlı. Başka bir takma ad seçmelisiniz!",
			}, nil
		}
		return nil, nil
	}
*/

func (s *urlService) GetById(id int) (*ent.Url, error) {
	return s.repo.GetById(id)
}
func (s *urlService) Delete(id int) (*dto.GeneralResponse[any], error) {

	err := s.repo.Delete(id)
	if err != nil {
		return &dto.GeneralResponse[any]{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	return &dto.GeneralResponse[any]{
		Status: true,
	}, nil
}
func (s *urlService) FindByShortCode(shortcode string) (*ent.Url, error) {

	url, err := s.repo.FindByShortCode(shortcode)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) Create(userID int, req dto.CreateUrlRequest, type_ string) (*dto.GeneralResponse[dto.UrlResponse], error) {
	// 1. URL formatını kontrol et
	u, err := netURL.ParseRequestURI(req.LongUrl) // req içindeki alan adının LongURL olduğunu varsayıyorum
	if err != nil {
		return nil, errors.New("geçersiz URL formatı")
	}

	// 2. Scheme (protokol) kontrolü (http veya https zorunluluğu)
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("URL mutlaka http veya https ile başlamalıdır")
	}

	// 3. Host (alan adı) kontrolü
	if u.Host == "" || !strings.Contains(u.Host, ".") {
		return nil, errors.New("geçersiz alan adı")
	}

	if req.Alias != "" {
		urll, err := s.repo.FindByAlias(req.Alias)
		if err == nil && urll != nil {
			return nil, errors.New("Takma ad kayıtlı. Başka bir takma ad seçmelisiniz!")
		}
	}

	urlResult, err := s.repo.Create(userID, utils.GenerateShortID(5), req)
	if err != nil {
		return nil, err
	}

	return &dto.GeneralResponse[dto.UrlResponse]{
		Status: urlResult.ShortCode != "",
		Data:   dto.ToUrlResponse(urlResult, type_),
	}, nil
}

func (s *urlService) History(userID int) ([]*dto.UrlResponse, error) {

	urll, err := s.repo.History(userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUrlResponseList(urll), nil
}

func (s *urlService) HistoryById(userID int, id int) (*dto.UrlResponse, error) {

	post, err := s.repo.HistoryById(userID, id)
	if err != nil {
		return nil, err
	}

	return dto.ToUrlResponse(post, "url"), nil
}

func (s *urlService) Update(req dto.UpdateUrlRequest) (*dto.GeneralResponse[any], error) {

	if req.Alias != "" {
		urll, err := s.repo.FindByAliasWithId(req.ID, req.Alias)
		if err == nil && urll != nil {
			return nil, errors.New("Takma ad kayıtlı. Başka bir takma ad seçmelisiniz!")
		}
	}

	post, err := s.repo.Update(req)
	if err != nil {
		return nil, err
	}

	return &dto.GeneralResponse[any]{
		Status: post.ShortCode != "",
	}, nil
}

func (s *urlService) Bulkcreate(userID int, req []dto.CreateUrlRequest) (*dto.GeneralResponse[dto.UrlResponse], error) {

	// 1. Liste içindeki her bir URL'yi kontrol et
	for _, item := range req {
		// URL'yi temizle ve parse et
		trimmedURL := strings.TrimSpace(item.LongUrl)
		u, err := netURL.ParseRequestURI(trimmedURL)

		// Temel format kontrolü
		if err != nil {
			return nil, errors.New("geçersiz URL formatı: " + item.LongUrl)
		}

		// Protokol kontrolü (http/https)
		if u.Scheme != "http" && u.Scheme != "https" {
			return nil, errors.New("URL mutlaka http veya https ile başlamalıdır: " + item.LongUrl)
		}

		// Host (alan adı) ve nokta kontrolü
		if u.Host == "" || !strings.Contains(u.Host, ".") {
			return nil, errors.New("geçersiz alan adı: " + item.LongUrl)
		}

		// Opsiyonel: ExpirationDate geçmişte mi kontrolü
		if !item.ExpirationDate.IsZero() && item.ExpirationDate.Before(time.Now()) {
			return nil, errors.New("son kullanma tarihi geçmiş olamaz: " + item.LongUrl)
		}

		if item.Alias != "" {
			urll, err := s.repo.FindByAlias(item.Alias)
			if err == nil && urll != nil {
				return nil, errors.New("Takma ad kayıtlı. Başka bir takma ad seçmelisiniz! Takma ad: " + item.Alias)
			}
		}

	}
	_, err := s.repo.CreateBulk(userID, req)
	if err != nil {
		return nil, err
	}

	return &dto.GeneralResponse[dto.UrlResponse]{
		Status:  true,
		Message: "Success",
	}, nil
}

func (s *urlService) Redirect(shortcodeOrAlias, password string) (*dto.GeneralResponse[dto.UrlResponse], error) {
	url, err := s.repo.FindShortCodeOrAlias(shortcodeOrAlias)
	if err != nil {
		return nil, err
	}

	// Şifre kontrolü
	if url.IsEncrypted && url.Password != password {
		return &dto.GeneralResponse[dto.UrlResponse]{
			Status:  false,
			Message: "Url is Encrypted",
			Code:    409,
		}, nil
	}

	// Süre kontrolü
	// 1. Önce ExpirationDate'in nil olup olmadığını kontrol ediyoruz (yani tarihin set edilip edilmediğini)
	if !url.ExpirationDate.IsZero() && time.Now().After(url.ExpirationDate) {
		return &dto.GeneralResponse[dto.UrlResponse]{
			Status:  false,
			Message: "Url is Expired",
			Code:    409,
		}, nil
	}

	// Veriyi hazırlıyoruz
	urlData := &dto.UrlResponse{LongUrl: url.LongURL, Id: url.ID}

	// Return tipine DİKKAT: Fonksiyon imzasıyla (any) aynı olmalı
	return &dto.GeneralResponse[dto.UrlResponse]{
		Status:  true,
		Message: "Success",
		Data:    urlData, // any tipinde kabul edilir
		Code:    200,
	}, nil
}
