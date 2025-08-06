package storage

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mrjvadi/BackendPanelVpn/cache"
	"github.com/mrjvadi/BackendPanelVpn/models"
	"github.com/mrjvadi/BackendPanelVpn/pkg"
	"github.com/mrjvadi/BackendPanelVpn/types"
	type_api "github.com/mrjvadi/BackendPanelVpn/types/type-api"
	type_redis "github.com/mrjvadi/BackendPanelVpn/types/type-redis"
	"time"
)

// Auth handles user authentication and session storage.
type Auth struct {
	db     *Database
	redis  *cache.RedisCache
	jwtKey []byte
}

// IAuth defines available auth operations.
type IAuth interface {
	// Login authenticates credentials, returns a JWT token and roles.
	Login(ctx context.Context, username, password string) types.BaseResponse[type_api.LoginResponse]

	// Validate checks if a session token is valid.
	Validate(ctx context.Context, token string) bool

	// Logout revokes a session token.
	Logout(ctx context.Context, token string) error

	// GetSessionByUsername retrieves the active session for a given username.
	GetSessionByUsername(ctx context.Context, username string) types.BaseResponse[type_api.LoginResponse]
}

// NewAuth constructs a new Auth with DB, Redis cache, and JWT secret.
func NewAuth(db *Database, redis *cache.RedisCache, jwtSecret string) IAuth {
	return &Auth{
		db:     db,
		redis:  redis,
		jwtKey: []byte(jwtSecret),
	}
}

// Login checks credentials, collects roles, generates JWT, stores it in Redis, and returns it.
func (a *Auth) Login(ctx context.Context, username, password string) types.BaseResponse[type_api.LoginResponse] {
	hashed := pkg.HashMD5(password)
	tx := a.db.Db.Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var roles []string

	// telegram id
	var telegramID int64

	// Check reseller role
	var reseller models.Reseller
	if err := tx.Where("username = ? AND password = ?", username, hashed).First(&reseller).Error; err == nil {
		roles = append(roles, "seller")
		if reseller.CanCreateSubReseller {
			roles = append(roles, "reseller")
		}
		telegramID = reseller.TelegramID
	}

	// Check admin role
	var admin models.Admin
	if reseller.ID == 0 { // Only check admin if not found as reseller
		if err := tx.Where("username = ? AND password = ?", username, hashed).First(&admin).Error; err == nil {
			roles = []string{"admin", "reseller", "seller"}
			telegramID = admin.TelegramID
		}
	}

	// If no roles found, unauthorized
	if len(roles) == 0 {
		tx.Rollback()
		return types.Error[type_api.LoginResponse](types.Unauthorized, types.LoginFailedMessage)
	}

	tx.Commit()
	return a.issueSession(ctx, username, roles, telegramID)
}

// issueSession generates a new JWT, revokes any existing session for the user, stores the new session in Redis, and returns the response.
func (a *Auth) issueSession(ctx context.Context, username string, roles []string, telegramID int64) types.BaseResponse[type_api.LoginResponse] {
	// ساخت JWT
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   username,
		"roles": roles,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := tokenObj.SignedString(a.jwtKey)
	if err != nil {
		panic(err)
	}

	// پاک کردن سشن قبلی (اگر ایندکس username داشتیم)
	if oldToken, _ := a.redis.GetUserIndex(ctx, username); oldToken != "" {
		a.redis.DeleteSessionData(ctx, oldToken)
	}

	// ذخیره‌ی سشن کامل و ایندکس کاربر → توکن
	session := &type_redis.SessionData{
		Token:      signed,
		Roles:      roles,
		TelegramID: telegramID,
		Username:   username,
	}
	a.redis.SetSessionData(ctx, session, 24*time.Hour)
	a.redis.SetUserIndex(ctx, username, signed, 24*time.Hour)

	return types.Ok(types.BaseResponse[type_api.LoginResponse]{
		Code:    types.Success,
		Message: types.LoginSuccessMessage,
		Data:    type_api.LoginResponse{Token: signed, Role: roles, TelegramID: telegramID},
	})
}
func (a *Auth) Validate(ctx context.Context, token string) bool {
	sess, err := a.redis.GetSessionData(ctx, token)
	return err == nil && sess != nil
}

func (a *Auth) GetSessionByUsername(ctx context.Context, username string) types.BaseResponse[type_api.LoginResponse] {
	token, _ := a.redis.GetUserIndex(ctx, username)
	if token == "" {
		return types.Error[type_api.LoginResponse](types.Unauthorized, types.LoginFailedMessage)
	}
	sess, err := a.redis.GetSessionData(ctx, token)
	if err != nil || sess == nil {
		return types.Error[type_api.LoginResponse](types.Unauthorized, types.SessionErrorMessageInvalidOrExpired)
	}
	return types.Ok(types.BaseResponse[type_api.LoginResponse]{
		Code:    types.Success,
		Message: types.LoginSuccessMessage,
		Data:    type_api.LoginResponse{Token: sess.Token, Role: sess.Roles, TelegramID: sess.TelegramID},
	})
}

func (a *Auth) Logout(ctx context.Context, token string) error {
	sess, err := a.redis.GetSessionData(ctx, token)
	if err != nil || sess == nil {
		return errors.New("session not found")
	}
	a.redis.DeleteSessionData(ctx, token)
	a.redis.DeleteUserIndex(ctx, sess.Username)
	return nil
}
