package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken verifies the token and returns the payload.
	VerifyToken(token string) (*Payload, error)
}
