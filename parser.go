package structparse

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	structTagParse     = "parse"
	structTagDelimiter = "delimiter"
)

type Config struct {
	Src           Source
	KeyFmt        KeyFmt
	Transformers  []Transformer
	IgnoreMissing bool
}

func Parse(cfg Config, dst interface{}) error {
	if cfg.Src == nil {
		return errors.New("structparse: source is missing")
	}
	if cfg.KeyFmt == nil {
		return errors.New("structparse: key formatter is missing")
	}
	return parse(cfg, dst, nil)
}

type ParseError []*FieldError

func (err ParseError) Error() string {
	if len(err) == 0 {
		return ""
	}
	all := make([]string, 0, len(err))
	for _, err := range err {
		all = append(all, err.Error())
	}
	return strings.Join(all, ", ")
}

func (err ParseError) Format(s fmt.State, verb rune) {
	if len(err) == 0 {
		return
	}
	if verb == 'v' && s.Flag('+') {
		for i, err := range err {
			if i > 0 {
				_, _ = fmt.Fprint(s, "\n")
			}
			_, _ = fmt.Fprint(s, err.Error())
		}
		return
	}
	if verb == 'v' || verb == 's' {
		_, _ = fmt.Fprint(s, err.Error())
	}
}

func parse(cfg Config, dst interface{}, parentKeys []string) error {
	v := reflect.Indirect(reflect.ValueOf(dst))
	if v.Kind() != reflect.Struct {
		return errors.New("structparse: dst is not a struct")
	}

	var retErr ParseError

	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		fieldValue := v.Field(i)

		err := parseField(cfg, fieldType, fieldValue, parentKeys)
		if errors.Is(err, ErrSourceKeyNotFound) && cfg.IgnoreMissing {
			continue
		}
		if err != nil {
			if as := ParseError(nil); errors.As(err, &as) {
				retErr = append(retErr, as...)
				continue
			}
			if as := (*FieldError)(nil); errors.As(err, &as) {
				retErr = append(retErr, as)
				continue
			}
			return err
		}
	}

	if len(retErr) == 0 {
		return nil
	}
	return retErr
}

type FieldError struct {
	Cause     error
	FieldName string
	KeyName   string
}

func (err *FieldError) Unwrap() error {
	return err.Cause
}

func (err *FieldError) Error() string {
	return fmt.Sprintf("[key: %s] [field: %s] %s", err.KeyName, err.FieldName, err.Cause)
}

func parseField(cfg Config, field reflect.StructField, fieldValue reflect.Value, parentKeys []string) error {
	name := field.Name
	if override, exists := field.Tag.Lookup(structTagParse); exists {
		name = override
	}

	// support skipping fields
	if name == "-" {
		return nil
	}

	// handle embedded structs
	if field.Anonymous {
		switch field.Type.Kind() {
		case reflect.Struct:
			if !fieldValue.CanAddr() || !fieldValue.Addr().CanInterface() {
				return nil
			}
			return parse(cfg, fieldValue.Addr().Interface(), parentKeys)
		case reflect.Ptr:
			if fieldValue.IsNil() {
				if !fieldValue.CanSet() {
					return nil
				}
				fieldValue.Set(reflect.New(field.Type.Elem()))
			}
			return parse(cfg, fieldValue.Elem().Addr().Interface(), parentKeys)
		default:
			return &FieldError{
				Cause:     fmt.Errorf("unsupported anonymus type %s", field.Type.Kind()),
				FieldName: field.Type.Name(),
			}
		}
	}

	if !fieldValue.CanSet() {
		return nil
	}

	parentKeys = append(parentKeys, name)

	// handle nested structs that do not have a custom parser
	switch field.Type.Kind() {
	case reflect.Struct:
		if _, ok := getCustomParser(field.Type); !ok {
			return parse(cfg, fieldValue.Addr().Interface(), parentKeys)
		}
	case reflect.Ptr:
		if fieldValue.Type().Elem().Kind() == reflect.Struct {
			if _, ok := getCustomParser(fieldValue.Type().Elem()); !ok {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(field.Type.Elem()))
				}
				return parse(cfg, fieldValue.Elem().Addr().Interface(), parentKeys)
			}
		}
	}

	key := cfg.KeyFmt.Format(parentKeys)
	value, err := cfg.Src.Get(key)
	for _, t := range cfg.Transformers {
		value, err = t.Transform(key, value, err, field.Tag)
	}
	if errors.Is(err, ErrTransformerSkipKey) {
		return nil
	}
	if err != nil {
		return &FieldError{Cause: err, FieldName: field.Name, KeyName: key}
	}

	err = assignValue(fieldValue, value, field.Tag)
	if err != nil {
		return &FieldError{Cause: err, FieldName: field.Name, KeyName: key}
	}

	return nil
}

type AssignError struct {
	Cause error
	Msg   string
}

func (err *AssignError) Unwrap() error {
	return err.Cause
}

func (err *AssignError) Error() string {
	return fmt.Sprintf("%s: %s", err.Msg, err.Cause)
}

func assignValue(dst reflect.Value, src string, tag reflect.StructTag) error {
	typ := dst.Type()

	if p, ok := getCustomParser(typ); ok {
		val, err := p(src, tag)
		if err != nil {
			return &AssignError{
				Cause: err,
				Msg:   fmt.Sprintf("cannot parse as custom type (%s %s)", typ.PkgPath(), typ.Name()),
			}
		}
		dst.Set(reflect.ValueOf(val))
		return nil
	}

	switch typ.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(src)
		if err != nil {
			return &AssignError{Cause: err, Msg: "cannot parse as bool"}
		}
		dst.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(src, 0, typ.Bits())
		if err != nil {
			return &AssignError{Cause: err, Msg: "cannot parse as int"}
		}
		dst.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(src, 0, typ.Bits())
		if err != nil {
			return &AssignError{Cause: err, Msg: "cannot parse as uint"}
		}
		dst.SetUint(v)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(src, typ.Bits())
		if err != nil {
			return &AssignError{Cause: err, Msg: "cannot parse as float"}
		}
		dst.SetFloat(f)

	case reflect.Map:
		values, err := url.ParseQuery(src)
		if err != nil {
			return &AssignError{Cause: err, Msg: "cannot parse as map"}
		}

		dst.Set(reflect.MakeMap(typ))

		for k, v := range values {
			if len(v) == 0 {
				continue
			}
			keyVal := reflect.New(typ.Key())
			err := assignValue(keyVal, k, tag)
			if err != nil {
				return err
			}
			elemVal := reflect.New(typ.Elem())
			err = assignValue(elemVal, v[0], tag)
			if err != nil {
				return err
			}
			dst.SetMapIndex(keyVal.Elem(), elemVal.Elem())
		}

	case reflect.Ptr:
		if dst.IsNil() {
			dst.Set(reflect.New(typ.Elem()))
		}
		return assignValue(dst.Elem(), src, tag)

	case reflect.Slice:
		if len(src) == 0 {
			return nil
		}
		if typ.Elem().Kind() == reflect.Uint8 {
			// handle []byte slices as string
			dst.SetBytes([]byte(src))
		} else {
			delimiter := tag.Get(structTagDelimiter)
			if len(delimiter) == 0 {
				delimiter = ","
			}
			parts := strings.Split(src, delimiter)
			sl := reflect.MakeSlice(typ, len(parts), len(parts))
			for i, p := range parts {
				err := assignValue(sl.Index(i), p, tag)
				if err != nil {
					return err
				}
			}
			dst.Set(sl)
		}

	case reflect.String:
		dst.SetString(src)

	default:
		return fmt.Errorf("unsupported dst type %s", typ.Kind())
	}

	return nil
}
