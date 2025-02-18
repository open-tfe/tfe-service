package constants

import "time"

// HTTP related constants
const (
	DefaultTimeout = 10 * time.Second
	APIVersionPath = "/api/v2"
)

// Context key constants
type ContextKey string

const (
	UserEmailKey ContextKey = "userEmail"
	UserTokenKey ContextKey = "userToken"
)
