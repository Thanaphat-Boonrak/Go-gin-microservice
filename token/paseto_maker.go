package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	paseto       paseto.Parser
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	maker := &PasetoMaker{
		symmetricKey: key,
		paseto:       paseto.NewParser(),
	}

	return maker, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := paseto.NewToken()

	token.SetString("id", payload.ID.String())
	token.SetString("username", payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpireAt)

	signed := token.V4Encrypt(p.symmetricKey, nil)

	return signed, nil
}
func (p *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := p.paseto.ParseV4Local(p.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	username, err := token.GetString("username")
	if err != nil {
		return nil, fmt.Errorf("invalid token: missing username claim")
	}

	idStr, err := token.GetString("id")
	if err != nil {
		return nil, fmt.Errorf("invalid token: missing id claim")
	}

	tokenID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token: invalid id format")
	}

	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("invalid token: missing issued_at")
	}

	expiredAt, err := token.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("invalid token: missing expiration")
	}

	payload := &Payload{
		ID:       tokenID,
		Username: username,
		IssuedAt: issuedAt,
		ExpireAt: expiredAt,
	}

	return payload, nil
}
