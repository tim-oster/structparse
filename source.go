package structparse

import (
	"errors"
	"net/url"
	"os"
)

var (
	ErrSourceKeyNotFound = errors.New("key not found")
)

type Source interface {
	Get(key string) (string, error)
}

type SourceFunc func(key string) (string, error)

func (src SourceFunc) Get(key string) (string, error) {
	return src(key)
}

func SourceEnv() Source {
	return SourceFunc(func(key string) (string, error) {
		if val, exists := os.LookupEnv(key); exists {
			return val, nil
		}
		return "", ErrSourceKeyNotFound
	})
}

type SourceMap map[string]string

func (src SourceMap) Get(key string) (string, error) {
	if val, exists := src[key]; exists {
		return val, nil
	}
	return "", ErrSourceKeyNotFound
}

func SourceUrl(values url.Values) Source {
	m := make(SourceMap)
	for k, v := range values {
		if len(v) == 0 {
			continue
		}
		m[k] = v[0]
	}
	return m
}

func SourceNil() Source {
	return SourceFunc(func(key string) (string, error) {
		return "", ErrSourceKeyNotFound
	})
}
