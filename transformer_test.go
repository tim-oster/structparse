package structparse

import (
	"reflect"
	"testing"
)

func TestTransformerDefaultValue(t *testing.T) {
	var dummy struct {
		WithoutDefault string
		WithDefault    string `default:"default_value"`
	}

	cfg := Config{
		Src:          SourceMap{},
		KeyFmt:       KeyFmtKebab(),
		Transformers: []Transformer{TransformerDefaultValue()},
	}
	err := Parse(cfg, &dummy)
	assertEqual(t, "[key: without-default] [field: WithoutDefault] key not found", err)
	assertEqual(t, "", dummy.WithoutDefault)
	assertEqual(t, "default_value", dummy.WithDefault)
}

func TestErrTransformerSkipKey(t *testing.T) {
	var dummy struct {
		Missing string
	}

	cfg := Config{
		Src:    SourceMap{},
		KeyFmt: KeyFmtKebab(),
		Transformers: []Transformer{TransformFunc(func(key, srcValue string, srcErr error, tag reflect.StructTag) (string, error) {
			return "", ErrTransformerSkipKey
		})},
	}
	err := Parse(cfg, &dummy)
	assertEqual(t, nil, err)
	assertEqual(t, "", dummy.Missing)
}
