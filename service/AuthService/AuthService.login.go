package AuthService

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrjvadi/BackendPanelVpn/events"
	"github.com/mrjvadi/BackendPanelVpn/storage"
	"github.com/mrjvadi/BackendPanelVpn/types"
	type_api "github.com/mrjvadi/BackendPanelVpn/types/type-api"
	type_event "github.com/mrjvadi/BackendPanelVpn/types/type-event"
	"net/http"
	"time"
)

type Auth struct {
	storage *storage.Store
	event   *events.Bus
}

type IAuth interface {
	Login(ctx *gin.Context, username, password string) types.BaseResponse[type_api.LoginResponse]
	Logout(token string) types.BaseResponse[any]
	//RefreshToken(token string) (string, error)
}

func NewAuth(event *events.Bus, storage *storage.Store) IAuth {
	return &Auth{
		storage: storage,
		event:   event,
	}
}

func (a *Auth) Login(ctx *gin.Context, username, password string) types.BaseResponse[type_api.LoginResponse] {
	ctxs, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	active := a.storage.Auth.GetSessionByUsername(ctxs, username)

	if active.IsSuccess() {
		// publish login event
		requestID := uuid.New().String()
		a.event.Publish(events.Event{
			Name: "any:login",
			Payload: type_event.EventBotLogin{
				Username:   username,
				RequestID:  requestID,
				LoginTime:  time.Now().Format(time.RFC3339) + " UTC",
				IPAddress:  ctx.ClientIP(),
				DeviceInfo: ctx.GetHeader("User-Agent"),
				Code:       http.StatusOK,
				TelegramID: active.Data.TelegramID,
			},
		})
		return active
	}

	active = a.storage.Auth.Login(ctxs, username, password)

	if active.IsSuccess() {
		// Publish login event
		requestID := uuid.New().String()
		a.event.Publish(events.Event{
			Name: "any:login",
			Payload: type_event.EventBotLogin{
				Username:   username,
				RequestID:  requestID,
				LoginTime:  time.Now().Format(time.RFC3339) + " UTC",
				IPAddress:  ctx.ClientIP(),
				DeviceInfo: ctx.GetHeader("User-Agent"),
				Code:       http.StatusOK,
				TelegramID: active.Data.TelegramID,
			},
		})
	}

	return active

}

func (a *Auth) Logout(token string) types.BaseResponse[any] {
	ctxs, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.storage.Auth.Logout(ctxs, token)

	if err != nil {
		return types.Error[any](types.Unauthorized, types.LoginFailedMessage)
	}

	return types.Ok[any](types.BaseResponse[any]{
		Code:    types.Success,
		Message: types.LogoutSuccessMessage,
		Data:    nil,
	})
}

//
//func (a *Auth) RefreshToken(token string) types.BaseResponse[type_api.LoginResponse] {
//	ctxs, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	active := a.storage.Auth.RefreshToken(ctxs, token)
//
//	if active.IsSuccess() {
//		// Publish refresh token event
//		requestID := uuid.New().String()
//		a.event.Publish(events.Event{
//			Name: "any:refresh_token",
//			Payload: type_event.EventBotLogin{
//				Username:   active.Data.Username,
//				RequestID:  requestID,
//				LoginTime:  time.Now().Format(time.RFC3339) + " UTC",
//				IPAddress:  "",
//				DeviceInfo: "",
//				Code:       http.StatusOK,
//				TelegramID: active.Data.TelegramID,
//			},
//			RequestID: requestID,
//		})
//	}
//
//	return active
//}
