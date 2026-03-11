package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

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

type UrlRepository interface {
	Save(sourceUrl string, shortenedUrl string) error
}

type UrlService struct {
	repo UrlRepository
}

func NewUrlService(repo UrlRepository) *UrlService {
	return &UrlService{
		repo: repo,
	}
}

func (s *UrlService) Save(sourceUrl string) (string, error) {
	if !(strings.HasPrefix(sourceUrl, "http://") || strings.HasPrefix(sourceUrl, "https://")) {
		return "", fmt.Errorf("not a link")
	}

	newPath := genRandomString(5)

	err := s.repo.Save(sourceUrl, newPath)

	if err != nil {
		if errors.Is(err, repository.ErrShortenedUrlAlreadyExists) {
			fmt.Println("recursive")
			return s.Save(sourceUrl)
		}
		return "", fmt.Errorf("can't save new url to storage: %w", err)
	}

	return newPath, nil
}
