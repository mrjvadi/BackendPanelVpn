// Filename: storage/store.go
package storage

import "github.com/mrjvadi/BackendPanelVpn/cache"

type Store struct {
	Auth IAuth
}

func NewStore(db *Database, ca *cache.RedisCache, jwtToken string) *Store {
	return &Store{
		Auth: NewAuth(db, ca, jwtToken),
	}
}
