package models

import "time"

type SessionData struct {
	ClientId     string    `json:"client_id"`
	Token        string    `json:"token"`
	LoginTime    time.Time `json:"login_time"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	DeviceType   string    `json:"device_type"`
	LastActivity time.Time `json:"last_activity"`
}
