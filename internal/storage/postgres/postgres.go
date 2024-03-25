package postgres

import (
	"errors"
	"github.com/brianvoe/gofakeit"
)

type Storage struct {
	urlMap map[string]string
}

func NewStorage() (*Storage, error) {
	return &Storage{
		urlMap: make(map[string]string),
	}, nil
}

func CreateUrl(url string, s *Storage) (string, error) {
	id := gofakeit.UUID()
	s.urlMap[id] = url

	return id, nil
}

func GetUrl(url string, s *Storage) (string, error) {
	fUrl, ok := s.urlMap[url]
	if !ok {
		return "", errors.New("key does not exists")
	}

	return fUrl, nil
}
