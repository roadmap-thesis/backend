package auth

import (
	"context"
	"time"

	"github.com/HotPotatoC/roadmap_gen/pkg/commonerrors"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthCtxKey = "identity"
)

type Payload struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(id int, email string, expiresIn time.Duration) *Payload {
	return &Payload{
		ID:        id,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(expiresIn),
	}
}

func NewPayloadFromClaims(claims jwt.MapClaims) *Payload {
	iat := int64(claims["iat"].(float64))
	exp := int64(claims["exp"].(float64))
	return &Payload{
		ID:        int(claims["id"].(float64)),
		Email:     claims["email"].(string),
		IssuedAt:  time.Unix(iat, 0),
		ExpiresAt: time.Unix(exp, 0),
	}
}

// FromContext extracts the auth payload
func FromContext(ctx context.Context) *Payload {
	return ctx.Value(AuthCtxKey).(*Payload)
}

func (p *Payload) Claims() jwt.Claims {
	return jwt.MapClaims{
		"id":    p.ID,
		"email": p.Email,
		"iat":   p.IssuedAt.Unix(),
		"exp":   p.ExpiresAt.Unix(),
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return commonerrors.Unauthorized()
	}

	return nil
}
