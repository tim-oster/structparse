package structparse

import (
	"fmt"
	"reflect"
	"sync"
)

type Parser func(value string, tag reflect.StructTag) (interface{}, error)

var (
	customParsersMu sync.Mutex
	customParsers   = make(map[string]Parser)
)

func RegisterCustomParser(typ interface{}, fn Parser) {
	customParsersMu.Lock()
	defer customParsersMu.Unlock()

	name := getCustomParserName(reflect.TypeOf(typ))
	if _, ok := customParsers[name]; ok {
		panic(fmt.Sprintf("custom parser already registered for: %s", name))
	}
	customParsers[name] = fn
}

func getCustomParserName(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.PkgPath() + "/" + typ.Name()
}

func getCustomParser(typ reflect.Type) (Parser, bool) {
	customParsersMu.Lock()
	defer customParsersMu.Unlock()

	parser, ok := customParsers[getCustomParserName(typ)]
	return parser, ok
}
