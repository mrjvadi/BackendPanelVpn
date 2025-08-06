package main

import (
	"github.com/mrjvadi/BackendPanelVpn/api"
	"github.com/mrjvadi/BackendPanelVpn/bot"
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/config"
	"github.com/mrjvadi/BackendPanelVpn/events"
	"github.com/mrjvadi/BackendPanelVpn/service"
	"github.com/mrjvadi/BackendPanelVpn/storage"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Cache
	ca, err := cache.NewRedisCache(cfg.Redis)

	if err != nil {
		panic("FATAL: Redis connection failed: " + err.Error())
	}

	// Initialize database connection
	condb := storage.ConnectPostgres(cfg.Postgres)

	// Initialize Database
	db := storage.NewStore(condb, ca, cfg.API.JwtToken)

	//condb.Db.Create(&models.Reseller{
	//	ParentResellerID:     nil,
	//	Username:             "reseller",
	//	Password:             pkg.HashMD5("reseller"), // This should be hashed in production
	//	TelegramID:           cfg.Bot.Admin,
	//	APIToken:             pkg.RandomString(32),
	//	CustomName:           pkg.RandomString(4),
	//	TotalQuota:           pkg.GB(20),
	//	Status:               "active",
	//	CanCreateSubReseller: false,
	//})

	// Initialize Event Bus
	bus := events.NewBus()

	// bot initialization

	bot := bot.New(bus, cfg.Bot)

	// Start the bot
	go bot.Start()

	// Service initialization

	siv := service.NewService(bus, db, ca)

	ap := api.NewApi(siv, &cfg.API, ca)
	ap.Init()

}
