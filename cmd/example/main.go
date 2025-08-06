// Filename: cmd/main.go
package main

import (
	"fmt"
	"github.com/mrjvadi/BackendPanelVpn/api/router"
	"log"
	"time"

	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/engine"
	"github.com/mrjvadi/BackendPanelVpn/events"
	"github.com/mrjvadi/BackendPanelVpn/models"
	"github.com/mrjvadi/BackendPanelVpn/service"
	"github.com/mrjvadi/BackendPanelVpn/storage"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// --- 1. Setup Database ---
	db, err := gorm.Open(sqlite.Open("final_project.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: DB connection failed: %v", err)
	}
	db.AutoMigrate(&models.User{})
	fmt.Println("✅ Database connected & migrated.")

	// --- 2. Setup Redis Cache (Low-level) ---
	redisStore, err := cache.NewRedisCache("localhost:6380")
	if err != nil {
		log.Fatalf("FATAL: Redis connection failed: %v", err)
	}
	fmt.Println("✅ Low-level Cache (Redis) connected.")

	// --- 3. Create a Typed Cache Wrapper for Users ---
	userTypedCache := cache.NewTypedCache[models.User](redisStore)
	fmt.Println("✅ Type-safe cache for Users created.")

	// --- 4. Setup Core Components ---
	bus := events.NewBus()
	appEngine := engine.New(bus)
	appEngine.Start()
	fmt.Println("✅ Event Bus & Engine started.")

	// --- 5. Setup Repositories and Services ---
	store := storage.NewStore(db)
	userService := service.NewUserService(store.User, userTypedCache, bus)
	jobService := service.NewJobService(bus)
	fmt.Println("✅ Store & Services initialized.")

	// --- 6. Setup API Router ---
	router := router.NewRouter(userService, redisStore)
	fmt.Println("✅ API router configured.")

	// --- 7. Start Server and Demonstrate Job ---
	go func() {
		fmt.Println("🚀 Starting API server on http://localhost:8080")
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("FATAL: API server failed: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("\n--- Demonstrating Async Job Processing ---")

	response, err := jobService.RequestJob("process_invoices", 5*time.Second)
	if err != nil {
		log.Printf("ERROR: Job request failed: %v", err)
	} else {
		fmt.Printf("✅ SUCCESS: Main app received callback: '%s'\n", response.Payload)
	}

	select {}
}
