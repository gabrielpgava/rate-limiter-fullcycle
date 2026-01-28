package models

import "time"

type IPstate struct {
	Count       int
	WindowStart time.Time
	BannedUntil time.Time
}
