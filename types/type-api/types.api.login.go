package type_api

// LoginRequest مدل درخواست لاگین
// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse مدل دادهٔ برگشتی پس از لاگین
// swagger:model LoginResponse
type LoginResponse struct {
	Token      string   `json:"token" binding:"required"`
	Role       []string `json:"role" binding:"required"` // نقش کاربر، مثلاً "admin" یا "reseller"
	TelegramID int64    `json:"-"`
}
