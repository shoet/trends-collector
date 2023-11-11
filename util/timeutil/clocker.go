package timeutil

import (
	"fmt"
	"time"

	"github.com/shoet/trends-collector/interfaces"
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

type FixedClocker struct{}

func (fc *FixedClocker) Now() time.Time {
	return time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
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
