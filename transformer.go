package structparse

import (
	"errors"
	"reflect"
)

const (
	structTagDefault = "default"
)

var (
	ErrTransformerSkipKey = errors.New("skip key")
)

type Transformer interface {
	Transform(key, srcValue string, srcErr error, tag reflect.StructTag) (string, error)
}

type TransformFunc func(key, srcValue string, srcErr error, tag reflect.StructTag) (string, error)

func (t TransformFunc) Transform(key, srcValue string, srcErr error, tag reflect.StructTag) (string, error) {
	return t(key, srcValue, srcErr, tag)
}

func TransformerDefaultValue() Transformer {
	return TransformFunc(func(key, srcValue string, srcErr error, tag reflect.StructTag) (string, error) {
		defaultValue, hasDefault := tag.Lookup(structTagDefault)
		if errors.Is(srcErr, ErrSourceKeyNotFound) && hasDefault {
			return defaultValue, nil
		}
		return srcValue, srcErr
	})
}
