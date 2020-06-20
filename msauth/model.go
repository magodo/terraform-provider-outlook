package msauth

import (
	"fmt"
	"time"
)

// Token is defined at: https://tools.ietf.org/html/rfc6749#section-5.1
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"` // "expires_in" is defined as RECOMMENDED, while in MSAUTH it is always returned, hence defined as `int`
	expiresOn    time.Time // The time when the access token is expired
	RefreshToken *string   `json:"refresh_token"`
	Scope        *string   `json:"scope"`
}

const (
	// General error code
	TokenErrorInvalidRequest       = "invalid_request"
	TokenErrorInvalidClient        = "invalid_client"
	TokenErrorInvalidGrant         = "invalid_grant"
	TokenErrorUnauthorizedClient   = "unauthorized_client"
	TokenErrorUnsupportedGrantType = "unsupported_grant_type"
	TokenErrorInvalidScope         = "invalid_scope"

	// Device flow specific error code
	TokenErrorAuthorizationPending = "authorization_pending"
	TokenErrorSlowDown             = "slow_down"
	TokenErrorAccessDenied         = "access_denied"
	TokenErrorExpiredToken         = "expired_token"
)

type TokenError struct {
	Error            string  `json:"error"`
	ErrorDescription *string `json:"error_description"`
	ErrorURI         *string `json:"error_uri"`
}

func (e TokenError) String() string {
	out := e.Error
	if e.ErrorDescription != nil {
		out = fmt.Sprintf("%s: %s", out, *e.ErrorDescription)
	}
	if e.ErrorURI != nil {
		out = fmt.Sprintf("%s (uri: %s)", out, *e.ErrorURI)
	}
	return out
}

type DeviceAuthorization struct {
	DeviceCode              string  `json:"device_code"`
	UserCode                string  `json:"user_code"`
	VerificationURI         string  `json:"verification_uri"`
	VerificationURIComplete *string `json:"verification_uri_complete"`
	ExpiresIn               int     `json:"expires_in"`
	Interval                *int    `json:"interval"`
}
