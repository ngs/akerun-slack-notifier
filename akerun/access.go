package akerun

import "time"

// Akerun .
type Akerun struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// User .
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// Access .
type Access struct {
	ID         int       `json:"id"`
	Action     string    `json:"action"`
	DeviceName string    `json:"device_name"`
	DeviceType string    `json:"device_type"`
	AccessedAt time.Time `json:"accessed_at"`
	Akerun     Akerun    `json:"akerun"`
	User       *User     `json:"user"`
}

// AccessesResponse .
type AccessesResponse struct {
	Accesses []Access `json:"accesses"`
}
