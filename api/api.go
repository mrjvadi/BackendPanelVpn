package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mrjvadi/BackendPanelVpn/api/router"
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"github.com/mrjvadi/BackendPanelVpn/service"
)

type Api struct {
	ginApi   *gin.Engine
	services *service.Service
	cfg      *config.APIConfig
}

type IApi interface {
	Init() string
}

func NewApi(srv *service.Service, cfg *config.APIConfig, redis *cache.RedisCache) IApi {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	// You can add more middleware here if needed
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-session-id")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content for preflight requests
			return
		}
		c.Next()
	})

	routers := router.NewRouter(r, cfg, srv, redis)

	if err := routers.Init(); err != nil {
		panic("Failed to initialize routers: " + err.Error())
	}

	return &Api{
		ginApi:   r,
		services: srv,
		cfg:      cfg,
	}

}

func (a *Api) Init() string {

	return a.ginApi.Run(fmt.Sprintf(":%s", a.cfg.Port)).Error()
}
