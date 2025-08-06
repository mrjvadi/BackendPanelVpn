package models

// Session represents an authenticated user session stored in cache.
type Session struct {
	Id   uint   `json:"id"`
	Type string `json:"type"`
}
