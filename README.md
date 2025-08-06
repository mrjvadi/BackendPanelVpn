# github.com/mrjvadi/BackendPanelVpn

## Application Architecture Overview

We have a modular application consisting of four main components:

1. **Web Service**
2. **Engine**
3. **API**
4. **Bot**

All of these components communicate via a shared **Service** layer. Both the Bot, the Engine, and the external API call into this Service.

---

## 1. Event-Driven Communication

To coordinate interactions between Service, Engine, API and Bot, we use an event-driven pattern:

1. **Publishing an Event**  
   Any component (Service or Engine) can publish an event.
2. **Request to the Engine**  
   The published event triggers a request to the Engine.  
   Each request carries a **UUID** that serves as a callback identifier.
3. **Callback & Response**  
   Once the Engine processes the request, it looks up the event data using that UUID and returns the result via the callback.
4. **Bidirectional Flow**  
   Events may flow:
   * Service → Engine
   * Engine → Service
   * Engine → Bot  
   In every case, the UUID lets the caller retrieve the response.

---

## 2. “Grokking” Operations Within the Engine

Inside the Engine we perform advanced "grokking" operations (e.g. pattern matching, data enrichment, analytics). These computations are encapsulated behind a clean API that the Service and Bot can invoke.

---

## 3. API Layer (Gin Framework)

Our HTTP API is implemented in Go using [Gin](https://github.com/gin-gonic/gin). The structure is:

* **Handler** – each endpoint has a dedicated handler function.
* **Router** – routes are grouped under a `gin.RouterGroup`.
* **Middleware** – session authentication is enforced via Redis-backed middleware.

Strict request/response types in Go (e.g. `dto.RequestType`, `dto.ResponseType`) ensure requests are validated and responses follow a uniform JSON structure.

### 3.1 Caching Strategy

To optimize read performance the API uses a read-through cache:

1. **Cache-First Lookups** – each read checks Redis first.
2. **Cache Miss Handling** – missing values are loaded from the database and then stored back into Redis.

### 3.2 Middleware: SessionAuthMiddleware

```go
func SessionAuthMiddleware(cache cache.CacheInterface) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1) Read the session ID from the HTTP header
        sessionID := c.GetHeader("x-session-id")

        // 2) Look up the session in Redis
        raw, err := cache.GetSession(context.Background(), sessionID)
        if err == redis.Nil {
            // Session is missing or expired
            c.AbortWithStatusJSON(http.StatusUnauthorized, dto.BaseResponse[any]{
                Status:  dto.ErrorStatus,
                Message: dto.SessionErrorMessageInvalidOrExpired,
            })
            return
        }
        if err != nil {
            // Internal error when accessing Redis
            c.AbortWithStatusJSON(http.StatusInternalServerError, dto.BaseResponse[any]{
                Status:  dto.ErrorStatus,
                Message: dto.SessionErrorMessageInternalError,
            })
            return
        }

        // 3) Decrypt the session type field
        decryptedType, _ := pkg.DecryptPassword(raw.Type, fmt.Sprintf("%v", raw.Id))
        raw.Type = decryptedType

        // 4) Store session info in the request context
        c.Set("session", raw)
        c.Set("sessionId", sessionID)

        c.Next()
    }
}
```

### 3.3 Registered Routes

```go
// Dashboard & Server Management
c.GET("/dashboard",      middleware.SessionAuthMiddleware(cache), h.DashboardHandler)
c.GET("/server",         middleware.SessionAuthMiddleware(cache), h.ServerHandler)
c.POST("/server",        middleware.SessionAuthMiddleware(cache), h.CreateServerHandler)

// Inbound Filter Management (Admins Only)
c.POST("/inbound/filter",   middleware.SessionAuthMiddleware(cache), h.FilterInboundAdminHandler)
c.DELETE("/inbound/filter", middleware.SessionAuthMiddleware(cache), h.RemoveInboundAdminHandler)

// URL Filter Management (Admins Only)
c.POST("/url/filter",      middleware.SessionAuthMiddleware(cache), h.FilterURLAdminHandler)
c.DELETE("/url/filter",    middleware.SessionAuthMiddleware(cache), h.RemoveURLAdminHandler)

// Reseller Management
c.GET("/reseller",         middleware.SessionAuthMiddleware(cache), h.ResellerHandler)

// Configuration & Tag Management
c.GET("/config",           middleware.SessionAuthMiddleware(cache), h.ConfigHandler)
c.POST("/config/save",     middleware.SessionAuthMiddleware(cache), h.ConfigSaveAdminHandler)
c.GET("/tag",              middleware.SessionAuthMiddleware(cache), h.TagHandler)
c.POST("/tag/save",        middleware.SessionAuthMiddleware(cache), h.TagSaveAdminHandler)
```

---

## 4. Configuration Management

At startup we load a central configuration object (via environment variables or a tool like Viper) that includes:

* Database connection details
* Redis connection details
* Telegram bot token
* Telegram channel ID for logs
* Administrator user IDs

This configuration is injected into both the Service and the Engine and is also accessible by the Bot.

---

## 5. Logging

A robust logging solution ensures:

1. **Persistent Storage** – all logs are written to the database for auditing.
2. **Real-Time Alerts** – critical errors are forwarded to the Bot which publishes them to a dedicated Telegram "logs" channel.

---

## 6. Documentation

* `README.md` always reflects the current system design and API reference.
* Whenever data models, API signatures, or workflows change, this documentation must be updated accordingly.

---

## 7. API Inputs & Outputs on the Web

An interactive web interface (such as Swagger UI or Redoc) displays all API endpoints, shows example request/response schemas, and allows real-time "Try It Out" testing against the development server.

---

## Summary

This architecture ensures:

* **Modularity & Decoupling** via event-driven communication
* **Performance** with cache-first lookups
* **Security** through session middleware and typed requests
* **Reliability** via persistent logs and real-time alerts
* **Maintainability** through documentation and configuration best practices

