package interfaces

import (
	"time"
)
type Clocker interface {
	Now() time.Time
}
