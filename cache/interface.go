// Filename: cache/interface.go
package cache

import (
	"context"
	"github.com/mrjvadi/BackendPanelVpn/models"
)

// CacheInterface defines cache operations for sessions.
type CacheInterface interface {
	// GetSession retrieves raw session data by token
	GetSession(ctx context.Context, token string) (models.Session, error)
	// GetUserSession retrieves the active JWT for a given username
	GetUserSession(ctx context.Context, username string) (string, error)
}
