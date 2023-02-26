package config

import "time"

const (
	DatabaseQueryTimeLayout = `'YYYY-MM-DD"T"HH24:MI:SS"."MS"Z"TZ'`
	// DatabaseTimeLayout
	DatabaseTimeLayout string = time.RFC3339
	// AccessTokenExpiresInTime ...
	AccessTokenExpiresInTime time.Duration = 1 * 24 * 60 * time.Minute
	// RefreshTokenExpiresInTime ...
	RefreshTokenExpiresInTime time.Duration = 30 * 24 * 60 * time.Minute
	// RedisCacheTTL ...
	RedisCacheTTL time.Duration = 1 * 24 * 60 * time.Minute
)

var (
	SigningKey = []byte("FfLbN7pIEYe8@!EqrttOLiwa(H8)7Ddo")
)
