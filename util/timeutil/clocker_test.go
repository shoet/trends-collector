package timeutil

import (
	"testing"
)

func Test_NowFormat(t *testing.T) {
	clocker := FixedClocker{}

	s := NowFormatRFC3339(&clocker)
	want := "2023-01-01T00:00:00+09:00"
	if s != want {
		t.Errorf("failed NowFormatRFC3339(): got %s, want %s", s, want)
	}

	s = NowFormatYYYYMMDD(&clocker)
	want = "20230101"
	if s != want {
		t.Errorf("failed NowFormatYYYYMMDD(): got %s, want %s", s, want)
	}

	s = NowFormatYYYYMMDDHHMMSS(&clocker)
	want = "20230101000000"
	if s != want {
		t.Errorf("failed NowFormatYYYYMMDDHHMMSS(): got %s, want %s", s, want)
	}

}
