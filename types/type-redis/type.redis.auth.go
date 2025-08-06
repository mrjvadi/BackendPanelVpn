package type_redis

type SessionData struct {
	Token      string   `json:"token"`
	Roles      []string `json:"roles"`
	TelegramID int64    `json:"telegramId"`
	Username   string   `json:"username"`
}
