package timeutil

import (
	"fmt"
	"time"

	"github.com/shoet/trends-collector/interfaces"
	_ "time/tzdata"
)

type RealClocker struct{}

func NewRealClocker() (*RealClocker, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed time.LoadLocation: %w", err)
	}
	time.Local = loc
	return &RealClocker{}, nil
}

func (rc *RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct {
	year   int
	month  time.Month
	day    int
	hour   int
	minute int
	second int
}

func NewFixedClocker(
	year int, month time.Month, day int, hour int, minute int, second int,
) (*FixedClocker, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed time.LoadLocation: %w", err)
	}
	time.Local = loc
	return &FixedClocker{
		year:   year,
		month:  month,
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	}, nil
}

func (fc *FixedClocker) Now() time.Time {
	return time.Date(fc.year, fc.month, fc.day, fc.hour, fc.minute, fc.second, 0, time.Local)
}

func NowFormatRFC3339(c interfaces.Clocker) string {
	return c.Now().Format(time.RFC3339)
}

func NowFormatYYYYMMDD(c interfaces.Clocker) string {
	return c.Now().Format("20060102")
}

func NowFormatYYYYMMDDHHMMSS(c interfaces.Clocker) string {
	return c.Now().Format("20060102150405")
}
