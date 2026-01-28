package models

import "time"

type TokenState struct {
	Count       int
	WindowStart time.Time
	BannedUntil time.Time
}
