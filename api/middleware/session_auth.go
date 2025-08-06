package middleware

import (
	"fmt"
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/types"
	type_api "github.com/mrjvadi/BackendPanelVpn/types/type-api"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

// JWTAuthMiddleware validates JWT from Authorization header,
// enforces single active JWT session per user, and injects session info into context.
func JWTAuthMiddleware(cache *cache.RedisCache, jwtSecret string) gin.HandlerFunc {
	key := []byte(jwtSecret)
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract Bearer token
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.Unauthorized,
				Message: types.SessionErrorMessageInvalidOrExpired,
			})
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		// Parse and validate JWT (includes exp check)
		tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.Unauthorized,
				Message: types.SessionErrorMessageInvalidOrExpired,
			})
			return
		}

		// Extract claims
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.Unauthorized,
				Message: types.SessionErrorMessageInvalidOrExpired,
			})
			return
		}

		// Subject (username)
		subVal, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.Unauthorized,
				Message: types.SessionErrorMessageInvalidOrExpired,
			})
			return
		}
		username := subVal

		// Roles as []string
		rRaw, ok := claims["roles"].([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.Unauthorized,
				Message: types.SessionErrorMessageInvalidOrExpired,
			})
			return
		}
		roles := make([]string, 0, len(rRaw))
		for _, r := range rRaw {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}

		// Optional: extract expiration
		expFloat, ok := claims["exp"].(float64)
		if ok {
			expTime := time.Unix(int64(expFloat), 0)
			c.Set("exp", expTime)
		}

		// Enforce single session: compare with stored JWT
		currentToken, err := cache.GetUserIndex(ctx, username)
		if err != nil {
			if err == redis.Nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
					Code:    types.Unauthorized,
					Message: types.SessionErrorMessageInvalidOrExpired,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.InternalServerError,
				Message: types.SessionErrorMessageInternalError,
			})
			return
		}
		if currentToken != tokenString {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.BaseResponse[type_api.BaseResponseError]{
				Code:    types.Unauthorized,
				Message: types.SessionErrorMessageInvalidOrExpired,
			})
			return
		}

		// Inject session data
		c.Set("token", tokenString)
		c.Set("username", username)
		c.Set("roles", roles)

		c.Next()
	}
}
