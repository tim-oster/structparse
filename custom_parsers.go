package structparse

import (
	"reflect"
	"time"
)

const (
	structTagLayout = "layout"
)

func init() {
	RegisterCustomParser(time.Duration(0), func(value string, tag reflect.StructTag) (interface{}, error) {
		v, err := time.ParseDuration(value)
		if err != nil {
			return nil, err
		}
		return v, nil
	})

	timeLayoutMapping := map[string]string{
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
	}

	RegisterCustomParser(time.Time{}, func(value string, tag reflect.StructTag) (interface{}, error) {
		layout := tag.Get(structTagLayout)
		if mapped, ok := timeLayoutMapping[layout]; ok {
			layout = mapped
		}
		if len(layout) == 0 {
			layout = time.RFC3339
		}
		v, err := time.Parse(layout, value)
		if err != nil {
			return nil, err
		}
		return v, nil
	})
}
