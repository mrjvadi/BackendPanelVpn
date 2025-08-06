// Filename: api/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"github.com/mrjvadi/BackendPanelVpn/service"
)

type Router struct {
	router *gin.Engine
	srv    *service.Service
	cfg    *config.APIConfig
	redis  *cache.RedisCache
}

type IRouter interface {
	Init() error
}

func NewRouter(router *gin.Engine, cfg *config.APIConfig, srv *service.Service, redis *cache.RedisCache) IRouter {
	return &Router{
		router: router,
		cfg:    cfg,
		srv:    srv,
		redis:  redis,
	}
}

func (r *Router) Init() error {

	api := r.router.Group("/api")
	r.initV1(api, r.cfg, r.srv, r.redis)

	return nil
}
