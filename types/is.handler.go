package types

import "net/http"

type CodeHandler int

func (c CodeHandler) Int() int {
	return int(c)
}

const (
	// Success
	Success CodeHandler = http.StatusOK // 200

	// 2xx – Success
	Created   CodeHandler = http.StatusCreated   // 201
	NoContent CodeHandler = http.StatusNoContent // 204

	// 4xx – Client Errors
	BadRequest   CodeHandler = http.StatusBadRequest   // 400
	Unauthorized CodeHandler = http.StatusUnauthorized // 401
	Forbidden    CodeHandler = http.StatusForbidden    // 403
	NotFound     CodeHandler = http.StatusNotFound     // 404
	Conflict     CodeHandler = http.StatusConflict     // 409

	// 5xx – Server Errors
	InternalServerError CodeHandler = http.StatusInternalServerError // 500
)

type MessageHandler string

func (m MessageHandler) String() string {
	return string(m)
}

const (
	SuccessMessage                      MessageHandler = "Success"
	LoginSuccessMessage                 MessageHandler = "Login successful"
	LoginFailedMessage                  MessageHandler = "Login failed, please check your username and password"
	InvalidRequestMessage               MessageHandler = "Invalid request data"
	UnauthorizedMessage                 MessageHandler = "Unauthorized access"
	InternalServerErrorMessage          MessageHandler = "Internal server error"
	NotFoundMessage                     MessageHandler = "Resource not found"
	BadRequestMessage                   MessageHandler = "Bad request"
	ForbiddenMessage                    MessageHandler = "Forbidden access"
	ConflictMessage                     MessageHandler = "Conflict occurred"
	CreatedMessage                      MessageHandler = "Resource created successfully"
	UpdatedMessage                      MessageHandler = "Resource updated successfully"
	DeletedMessage                      MessageHandler = "Resource deleted successfully"
	NoContentMessage                    MessageHandler = "No content available"
	SessionErrorMessageInvalidOrExpired MessageHandler = "Session is invalid or expired, please log in again"
	SessionErrorMessageInternalError    MessageHandler = "Internal server error while handling session"
	LogoutFailedMessage                 MessageHandler = "Logout failed, please try again later"
	LogoutSuccessMessage                MessageHandler = "Logout successful, you have been logged out"
)

type BaseResponse[T any] struct {
	Code    CodeHandler    `json:"code,omitempty" binding:"required" default:"200"`
	Message MessageHandler `json:"message"`
	Data    T              `json:"data,omitempty"`
}

func Ok[T any](response BaseResponse[T]) BaseResponse[T] {
	return response
}
func Error[T any](code CodeHandler, message MessageHandler) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    code,
		Message: message,
	}
}

func ErrorWithMessage[T any](message MessageHandler, data T) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    data,
	}
}

func ErrorWithCode[T any](code CodeHandler, data T) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    code,
		Message: MessageHandler(""),
		Data:    data,
	}
}

func ErrorWithCodeAndMessage[T any](code CodeHandler, message MessageHandler, data T) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func SuccessWithMessage[T any](message MessageHandler, data T) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    Success,
		Message: message,
		Data:    data,
	}
}
func SuccessWithData[T any](data T) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    Success,
		Message: SuccessMessage,
		Data:    data,
	}
}

func SuccessWithMessageAndData[T any](message MessageHandler, data T) BaseResponse[T] {
	return BaseResponse[T]{
		Code:    Success,
		Message: message,
		Data:    data,
	}
}
func (r BaseResponse[T]) IsSuccess() bool {
	code := int(r.Code)
	return code >= http.StatusOK && code < http.StatusMultipleChoices // 200 ≤ code < 300
}

func (r BaseResponse[T]) IsError() bool {
	return !r.IsSuccess()
}
