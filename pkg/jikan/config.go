package jikan

import "time"

const (
	MaxAllowedHitPerDay = 60 * 60 * 24
	MaxSafeHitPerDay    = 60 * 60 * 20
	Endpoint            = "https://api.jikan.moe/v4"
	CooldownDuration    = time.Second
	DefaultLimit        = 25
	Timeout             = 60
)
