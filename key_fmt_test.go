package structparse

import (
	"strings"
	"testing"
)

func TestKeyFmtFunc(t *testing.T) {
	result := KeyFmtFunc(func(keys []string) string {
		assertEqual(t, []string{"one", "two"}, keys)
		return "test"
	}).Format([]string{"one", "two"})
	assertEqual(t, "test", result)
}

func TestKeyFmtPrefix(t *testing.T) {
	result := KeyFmtPrefix("prefix-", KeyFmtFunc(func(keys []string) string {
		return "test"
	})).Format(nil)
	assertEqual(t, "prefix-test", result)
}

func TestKeyFmtJoin(t *testing.T) {
	result := KeyFmtJoin("-", nil).Format([]string{"one", "TWO"})
	assertEqual(t, "one-TWO", result)

	result = KeyFmtJoin("-", func(s string) string { return strings.ToUpper(s) }).Format([]string{"one", "two"})
	assertEqual(t, "ONE-TWO", result)
}

func TestKeyFmtEnv(t *testing.T) {
	result := KeyFmtEnv().Format([]string{"one", "TWO"})
	assertEqual(t, "ONE_TWO", result)
}

func TestKeyFmtKebab(t *testing.T) {
	result := KeyFmtKebab().Format([]string{"one", "TWO"})
	assertEqual(t, "one-two", result)
}
