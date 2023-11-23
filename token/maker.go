package token

import "time"

type Maker interface {
	/// CreateToken create token for specified username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	/// VerifyToken checks if the token is valid or invalid
	VerifyToken(token string) (*Payload, error)
}
