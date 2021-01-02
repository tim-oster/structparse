package decimal

import (
	"github.com/shopspring/decimal"
	"github.com/tim-oster/structparse"
	"reflect"
)

func init() {
	structparse.RegisterCustomParser(decimal.Decimal{}, func(value string, tag reflect.StructTag) (interface{}, error) {
		v, err := decimal.NewFromString(value)
		if err != nil {
			return nil, err
		}
		return v, nil
	})
}
