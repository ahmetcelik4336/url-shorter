package services

import (
	"api/ent"
	"api/internal/repositories"
	"errors"
	"net/url"
	dto "shared/models"
	"shared/utils"
	"strings"
	"time"
)

type UrlService interface {
	FindByShortCode(shortcode string) (*ent.Url, error)
	Create(userID int, req dto.CreateUrlRequest) (*dto.UrlResponse, error)
	History(userID int) ([]*dto.UrlResponse, error)
	Update(req dto.UpdateUrlRequest) (*dto.UrlResponse, error)
	Redirect(shortcodeOrAlias, password string) (*dto.GeneralResponse[dto.UrlResponse], error)
	Bulkcreate(userID int, req []dto.CreateUrlRequest) (*dto.GeneralResponse[dto.UrlResponse], error)
}

type urlService struct {
	repo repositories.UrlRepository
}

func NewUrlService(r repositories.UrlRepository) UrlService {
	return &urlService{repo: r}
}

func (s *urlService) FindByShortCode(shortcode string) (*ent.Url, error) {

	url, err := s.repo.FindByShortCode(shortcode)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) Create(userID int, req dto.CreateUrlRequest) (*dto.UrlResponse, error) {
	// 1. URL formatını kontrol et
	u, err := url.ParseRequestURI(req.LongUrl) // req içindeki alan adının LongURL olduğunu varsayıyorum
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
	post, err := s.repo.Create(userID, utils.GenerateShortID(5), req)
	if err != nil {
		return nil, err
	}

	return dto.ToUrlResponse(post), nil
}

func (s *urlService) History(userID int) ([]*dto.UrlResponse, error) {

	post, err := s.repo.History(userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUrlResponseList(post), nil
}

func (s *urlService) Update(req dto.UpdateUrlRequest) (*dto.UrlResponse, error) {

	post, err := s.repo.Update(req)
	if err != nil {
		return nil, err
	}

	return dto.ToUrlResponse(post), nil
}

func (s *urlService) Bulkcreate(userID int, req []dto.CreateUrlRequest) (*dto.GeneralResponse[dto.UrlResponse], error) {

	// 1. Liste içindeki her bir URL'yi kontrol et
	for _, item := range req {
		// URL'yi temizle ve parse et
		trimmedURL := strings.TrimSpace(item.LongUrl)
		u, err := url.ParseRequestURI(trimmedURL)

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
