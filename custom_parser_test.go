package structparse

import (
	"reflect"
	"testing"
)

func TestRegisterCustomParser(t *testing.T) {
	type dummy struct{}

	// register and test retrieval as pointer and non-pointer
	RegisterCustomParser((*dummy)(nil), func(value string, _ reflect.StructTag) (interface{}, error) {
		return value, nil
	})
	p, ok := getCustomParser(reflect.TypeOf(&dummy{}))
	if assertEqual(t, true, ok) {
		result, _ := p("1", "")
		assertEqual(t, "1", result)
	}
	p, ok = getCustomParser(reflect.TypeOf(dummy{}))
	if assertEqual(t, true, ok) {
		result, _ := p("2", "")
		assertEqual(t, "2", result)
	}

	// already registered
	func() {
		defer func() {
			p := recover()
			assertEqual(t, "custom parser already registered for: github.com/tim-oster/structparse/dummy", p)
		}()
		RegisterCustomParser((*dummy)(nil), nil)
	}()
}
