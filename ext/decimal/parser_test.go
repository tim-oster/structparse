package decimal

import (
	"github.com/shopspring/decimal"
	"github.com/tim-oster/structparse"
	"testing"
)

func TestCustomParserDecimalDecimal(t *testing.T) {
	var dummy struct {
		Decimal decimal.Decimal
	}

	cfg := structparse.Config{
		Src:    structparse.SourceMap{"decimal": "10.33"},
		KeyFmt: structparse.KeyFmtKebab(),
	}
	err := structparse.Parse(cfg, &dummy)
	if err != nil {
		t.Fatalf("no error expected: %s", err)
	}

	if got := dummy.Decimal.String(); got != "10.33" {
		t.Fatalf("expected 10.33 but got %s", got)
	}
}
