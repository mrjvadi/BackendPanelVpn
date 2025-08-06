package type_api

// BaseResponse
// swagger:model BaseResponse
type BaseResponse[T any] struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// BaseResponseError
// swagger:model BaseResponseError
type BaseResponseError struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}
