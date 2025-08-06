package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mrjvadi/BackendPanelVpn/api/handler"
	"github.com/mrjvadi/BackendPanelVpn/api/middleware"
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"github.com/mrjvadi/BackendPanelVpn/docs"
	_ "github.com/mrjvadi/BackendPanelVpn/docs" // Import the docs package to generate Swagger documentation
	"github.com/mrjvadi/BackendPanelVpn/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type V1 struct {
	r       *gin.RouterGroup
	srv     *service.Service
	cfg     *config.APIConfig
	handler *handler.Handler
	redis   *cache.RedisCache
}

func (r *Router) initV1(rg *gin.RouterGroup, cfg *config.APIConfig, srv *service.Service, redis *cache.RedisCache) {

	v1api := rg.Group("/v1")

	hand := handler.NewHandler(srv)

	v1 := V1{
		r:       v1api,
		cfg:     cfg,
		srv:     srv,
		handler: hand,
		redis:   redis,
	}

	v1.Auth()
	v1.Dashboard()
	v1.User()
	v1.Server()
	v1.Tag()
	v1.Seller()
	v1.Config()
	v1.Settings()
	v1.Analise()

	docs.SwaggerInfo.Host = fmt.Sprintf("192.168.0.100:%s", r.cfg.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1.r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func (v *V1) Auth() {

	auth := v.r.Group("/auth")

	auth.POST("/login", v.handler.Login)
	auth.POST("/logout", middleware.JWTAuthMiddleware(v.redis, v.cfg.JwtToken), v.handler.Logout)

}

func (v *V1) Dashboard() {
	dashboard := v.r.Group("/dashboard")

	dashboard.GET("")

}

func (v *V1) User() {
	user := v.r.Group("/user")

	user.GET("")

}

func (v *V1) Server() {
	server := v.r.Group("/server")

	server.GET("")
	server.POST("")
	server.PUT("")
	server.DELETE("")
}

func (v *V1) Tag() {
	tag := v.r.Group("/tag")

	tag.GET("")
	tag.POST("")
	tag.PUT("")
	tag.DELETE("")
}

func (v *V1) Seller() {
	seller := v.r.Group("/seller")

	seller.GET("")
	seller.POST("")
	seller.PUT("")
	seller.DELETE("")
}

func (v *V1) Config() {
	config := v.r.Group("/config")

	config.GET("")
	config.POST("")
	config.PUT("")
	config.DELETE("")
}

func (v *V1) Settings() {
	settings := v.r.Group("/settings")

	settings.GET("")
	settings.POST("")
	settings.PUT("")
	settings.DELETE("")
}

func (v *V1) Analise() {
	analise := v.r.Group("/analise")

	analise.GET("")
	analise.POST("")
	analise.PUT("")
	analise.DELETE("")
}
