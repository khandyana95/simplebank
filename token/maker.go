package token

import "time"

type Maker interface {
	CreateToken(string, time.Duration) (string, *Payload, error)
	ValidateToken(string) (*Payload, error)
}
