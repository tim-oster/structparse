package structparse

import (
	"fmt"
	"testing"
)

func TestParse_Errors(t *testing.T) {
	var dummy struct {
		Nested struct {
			FieldA bool
			FieldB int
		}
		FieldC  bool
		Missing string
	}

	cfg := Config{
		Src: SourceMap{
			"nested-field-a": "invalid",
			"nested-field-b": "invalid",
			"field-c":        "invalid",
		},
		KeyFmt: KeyFmtKebab(),
	}
	err := Parse(cfg, &dummy)

	expected := "[key: nested-field-a] [field: FieldA] cannot parse as bool: strconv.ParseBool: parsing \"invalid\": invalid syntax\n" +
		"[key: nested-field-b] [field: FieldB] cannot parse as int: strconv.ParseInt: parsing \"invalid\": invalid syntax\n" +
		"[key: field-c] [field: FieldC] cannot parse as bool: strconv.ParseBool: parsing \"invalid\": invalid syntax\n" +
		"[key: missing] [field: Missing] key not found"
	assertEqual(t, expected, fmt.Sprintf("%+v", err))
}

func TestParse_Missing(t *testing.T) {
	var dummy struct {
		Missing string
	}

	cfg := Config{
		Src:           SourceMap{},
		KeyFmt:        KeyFmtKebab(),
		IgnoreMissing: false,
	}
	err := Parse(cfg, &dummy)
	assertEqual(t, "[key: missing] [field: Missing] key not found", err)

	cfg.IgnoreMissing = true
	err = Parse(cfg, &dummy)
	assertEqual(t, nil, err)
}

func TestParse_Types(t *testing.T) {
	type Embedded1 struct {
		EmbeddedVal1 string
	}
	type Embedded2 struct {
		EmbeddedVal2 string
	}
	type privateEmbedded1 struct {
		PrivateVal string
	}
	type privateEmbedded2 struct {
		PrivateVal string
	}
	var dummy struct {
		Embedded1
		*Embedded2

		privateEmbedded1
		*privateEmbedded2

		unexported string
		Ignored    string `parse:"-"`

		Nested struct {
			Nested struct {
				InnerMost string
			}
		}
		NestedPtr *struct {
			InnerMost string
		}
		Nested2 struct {
			InnerMost string
		} `parse:"NestedRenamed"`

		NotYetRenamed string `parse:"Renamed"`

		Bool bool

		Int   int
		Int8  int8
		Int16 int16
		Int32 int32
		Int64 int64

		Uint   uint
		Uint8  uint8
		Uint16 uint16
		Uint32 uint32
		Uint64 uint64

		Float32 float32
		Float64 float64

		Map map[string]int

		PtrString *string

		ByteSlice         []byte
		IntSlice          []int
		IntSliceDelimiter []int `delimiter:";"`

		String string
	}

	cfg := Config{
		Src: SourceMap{
			"embedded-val1": "e1",
			"embedded-val2": "e2",

			"private-val": "should not be here",

			"unexported": "should not be here",
			"ignored":    "should not be here",

			"nested-nested-inner-most":  "n2",
			"nested-ptr-inner-most":     "np",
			"nested-renamed-inner-most": "nr",

			"renamed": "now",

			"bool": "true",

			"int":   "-1",
			"int8":  "-2",
			"int16": "-3",
			"int32": "-4",
			"int64": "-5",

			"uint":   "1",
			"uint8":  "2",
			"uint16": "3",
			"uint32": "4",
			"uint64": "5",

			"float32": "1.23",
			"float64": "2.34",

			"map": "k1=1&k2=2",

			"ptr-string": "pointed",

			"byte-slice":          "as string",
			"int-slice":           "1,2,3,4",
			"int-slice-delimiter": "5;6;7;8",

			"string": "test",
		},
		KeyFmt: KeyFmtKebab(),
	}
	err := Parse(cfg, &dummy)
	assertEqual(t, nil, err)

	assertEqual(t, "e1", dummy.EmbeddedVal1)
	assertEqual(t, "e1", dummy.Embedded1.EmbeddedVal1)
	assertEqual(t, "e2", dummy.EmbeddedVal2)
	assertEqual(t, "e2", dummy.Embedded2.EmbeddedVal2)

	assertEqual(t, "", dummy.privateEmbedded1.PrivateVal)
	assertEqual(t, nil, dummy.privateEmbedded2)

	assertEqual(t, "", dummy.unexported)
	assertEqual(t, "", dummy.Ignored)

	assertEqual(t, "n2", dummy.Nested.Nested.InnerMost)
	assertEqual(t, "np", dummy.NestedPtr.InnerMost)
	assertEqual(t, "nr", dummy.Nested2.InnerMost)

	assertEqual(t, "now", dummy.NotYetRenamed)

	assertEqual(t, true, dummy.Bool)

	assertEqual(t, -1, dummy.Int)
	assertEqual(t, -2, dummy.Int8)
	assertEqual(t, -3, dummy.Int16)
	assertEqual(t, -4, dummy.Int32)
	assertEqual(t, -5, dummy.Int64)

	assertEqual(t, 1, dummy.Uint)
	assertEqual(t, 2, dummy.Uint8)
	assertEqual(t, 3, dummy.Uint16)
	assertEqual(t, 4, dummy.Uint32)
	assertEqual(t, 5, dummy.Uint64)

	assertEqual(t, 1.23, dummy.Float32)
	assertEqual(t, 2.34, dummy.Float64)

	assertEqual(t, map[string]int{"k1": 1, "k2": 2}, dummy.Map)

	assertEqual(t, "pointed", *dummy.PtrString)

	assertEqual(t, []byte("as string"), dummy.ByteSlice)
	assertEqual(t, []int{1, 2, 3, 4}, dummy.IntSlice)
	assertEqual(t, []int{5, 6, 7, 8}, dummy.IntSliceDelimiter)

	assertEqual(t, "test", dummy.String)
}
