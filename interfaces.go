package checkpoint

import "errors"

type TokenVerifier interface {
	Verify(token, ipAddress string) (bool, error)
}

var (
	ErrLookupFailure = errors.New("unable to look up the status of the token provided")
	ErrServerConfig  = errors.New("the token response has one or more configuration-related errors")
)
