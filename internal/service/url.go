package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/HungryArthur/go-shortener/internal/config"
	"github.com/HungryArthur/go-shortener/internal/repository"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genRandomString(length int) string {
	str := make([]rune, length)
	for i := range length {
		str[i] = rune(alphabet[int(rand.Int31())%len(alphabet)])
	}
	return string(str)
}

type URLRepository interface {
	Save(sourceURL string, shortenedURL string) error
	Get(shortenedURL string) (string, error)
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) Save(sourceURL string) (string, error) {
	if !strings.HasPrefix(sourceURL, "http://") && !strings.HasPrefix(sourceURL, "https://") {
		return "", fmt.Errorf("not a link")
	}

	shortCode := genRandomString(5)

	err := s.repo.Save(sourceURL, shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrShortenedURLAlreadyExists) {
			fmt.Println("recursive")
			return s.Save(sourceURL)
		}
		return "", fmt.Errorf("can't save new url to storage: %w", err)
	}

	baseURL := config.FlagBaseShortenedURLAddr
	if baseURL == "" {
		return "http://localhost" + config.FlagRunAddr + "/" + shortCode, nil
	}

	baseURL = strings.TrimSuffix(baseURL, "/")
	return fmt.Sprintf("%s/%s", baseURL, shortCode), nil
}

func (s *URLService) Get(shortenedURL string) (string, error) {
	sourceURL, err := s.repo.Get(shortenedURL)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrShortenedURLDoesntExist):
			return "", ErrShortenedURLDoesntExist
		default:
			return "", fmt.Errorf("unhandled error from repository: %w, %w", err, ErrCantGetURL)
		}
	}
	return sourceURL, nil
}
