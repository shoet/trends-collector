package interfaces

import (
	"net/http"
	"time"
)

type Clocker interface {
	Now() time.Time
}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}
