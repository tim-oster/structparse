package structparse

import "strings"

type KeyFmt interface {
	Format(keys []string) string
}

type KeyFmtFunc func(keys []string) string

func (f KeyFmtFunc) Format(keys []string) string {
	return f(keys)
}

func KeyFmtPrefix(prefix string, inner KeyFmt) KeyFmt {
	return KeyFmtFunc(func(keys []string) string {
		key := inner.Format(keys)
		return prefix + key
	})
}

func KeyFmtJoin(separator string, mapFn func(string) string) KeyFmt {
	return KeyFmtFunc(func(keys []string) string {
		if mapFn != nil {
			for i, key := range keys {
				keys[i] = mapFn(key)
			}
		}
		return strings.Join(keys, separator)
	})
}

func KeyFmtEnv() KeyFmt {
	return KeyFmtJoin("_", CamelToUpperSnake)
}

func KeyFmtKebab() KeyFmt {
	return KeyFmtJoin("-", CamelToLowerKebab)
}
