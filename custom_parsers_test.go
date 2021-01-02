package structparse

import (
	"testing"
	"time"
)

func TestCustomParserTimeDuration(t *testing.T) {
	var dummy struct {
		Duration time.Duration
	}

	cfg := Config{
		Src:    SourceMap{"duration": "10m"},
		KeyFmt: KeyFmtKebab(),
	}
	err := Parse(cfg, &dummy)
	if err != nil {
		t.Fatalf("no error expected: %s", err)
	}

	assertEqual(t, 10*time.Minute, dummy.Duration)
}

func TestCustomParserTimeTime(t *testing.T) {
	var dummy struct {
		T1 time.Time // fallback layout: RFC3339
		T2 time.Time `layout:"Stamp"` // use layout mapping
		T3 time.Time `layout:"2006"`  // use custom layout
	}

	cfg := Config{
		Src: SourceMap{
			"t1": "2012-03-12T01:02:03Z",
			"t2": "Feb 10 11:12:13",
			"t3": "2021",
		},
		KeyFmt: KeyFmtKebab(),
	}
	err := Parse(cfg, &dummy)
	if err != nil {
		t.Fatalf("no error expected: %s", err)
	}

	assertEqual(t, "2012-03-12T01:02:03Z", dummy.T1.Format(time.RFC3339))
	assertEqual(t, "0000-02-10T11:12:13Z", dummy.T2.Format(time.RFC3339))
	assertEqual(t, "2021-01-01T00:00:00Z", dummy.T3.Format(time.RFC3339))
}
