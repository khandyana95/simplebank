package token

import "time"

type Maker interface {
	CreateToken(string, time.Duration) (string, error)
	ValidateToken(string) (*Payload, error)
}
