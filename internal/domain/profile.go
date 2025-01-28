package domain

import (
	"html"
	"time"
)

const (
	ProfileTable = "profiles"
)

type Profile struct {
	ID     int
	Name   string
	Avatar string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProfile(name, avatar string) *Profile {
	return &Profile{
		Name:      name,
		Avatar:    avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func GetDefaultAvatar(name string) string {
	return "https://hostedboringavatars.vercel.app/api/beam?colors=1DA1F2,14171A,657786,F5F8FA&name=" + html.EscapeString(name)
}
