package type_event

type EventBotLogin struct {
	Username   string `json:"username"`    // Username of the user logging in
	Password   string `json:"password"`    // Password of the user logging in
	RequestID  string `json:"request_id"`  // Unique identifier for the request
	LoginTime  string `json:"login_time"`  // Timestamp of the login event
	IPAddress  string `json:"ip_address"`  // IP address of the user logging in
	DeviceInfo string `json:"device_info"` // Information about the device used for login
	Code       int    `json:"code"`        // Optional code for two-factor authentication or similar purposes
	TelegramID int64  `json:"telegram_id"` // Optional Telegram ID for user identification in the bot
}
