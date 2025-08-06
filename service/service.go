// Filename: service/service.go
package service

import (
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/events"
	"github.com/mrjvadi/BackendPanelVpn/service/AuthService"
	"github.com/mrjvadi/BackendPanelVpn/storage"
)

type Service struct {
	events      *events.Bus
	db          *storage.Store
	cache       *cache.RedisCache
	AuthService AuthService.IAuth
}

func NewService(bus *events.Bus, db *storage.Store, redis *cache.RedisCache) *Service {

	auth := AuthService.NewAuth(bus, db)

	return &Service{
		events:      bus,
		db:          db,
		cache:       redis,
		AuthService: auth,
	}
}
