package util

import (
	"time"

	"github.com/shoet/trends-collector/interfaces"
)

type RealClocker struct{}

func (rc *RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (fc *FixedClocker) Now() time.Time {
	return time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
}

func NowFormatISO8601(c interfaces.Clocker) string {
	return c.Now().Format(time.RFC3339)
}
