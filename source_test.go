package structparse

import (
	"errors"
	"net/url"
	"os"
	"testing"
)

func testSrc(t *testing.T, src Source) {
	value, err := src.Get("key")
	assertEqual(t, "value", value)
	assertEqual(t, nil, err)

	value, err = src.Get("missing")
	assertEqual(t, "", value)
	if errors.Is(err, ErrSourceKeyNotFound) {
		assertEqual(t, "key not found", err)
	} else {
		t.Fatalf("expected error for missing env var but got '%s'", err)
	}
}

func TestSourceEnv(t *testing.T) {
	err := os.Setenv("key", "value")
	if err != nil {
		t.Fatalf("unexpeted error: %s", err)
	}
	testSrc(t, SourceEnv())
}

func TestSourceMap(t *testing.T) {
	testSrc(t, SourceMap{"key": "value"})
}

func TestSourceUrl(t *testing.T) {
	testSrc(t, SourceUrl(url.Values{"key": []string{"value"}}))
}

func TestSourceNil(t *testing.T) {
	_, err := SourceNil().Get("missing")
	assertEqual(t, ErrSourceKeyNotFound, err)
}
