package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.ExpireAt}, nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.IssuedAt}, nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.IssuedAt}, nil
}

func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (p *Payload) GetSubject() (string, error) {
	return p.Username, nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	token, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       token,
		Username: username,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}

	return payload, nil
}
