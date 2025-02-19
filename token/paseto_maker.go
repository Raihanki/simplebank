package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// PasetoMaker is a struct that holds the necessary components for PASETO token creation and validation.
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// NewPasetoMaker initializes a new PasetoMaker with a symmetric key.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, err
	}

	return &PasetoMaker{symmetricKey: key}, nil
}

// CreateToken creates a new PASETO token for the given username and duration.
func (pm *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	token := paseto.NewToken()
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)
	token.SetSubject(payload.Username)
	token.SetString("token_id", payload.ID.String())

	// Sign the token with the symmetric key
	signedToken := token.V4Encrypt(pm.symmetricKey, nil)
	return signedToken, nil
}

// VerifyToken verifies the PASETO token and extracts the payload.
func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()

	parsedToken, err := parser.ParseV4Local(pm.symmetricKey, token, nil)
	if err != nil {
		return nil, err
	}

	// Extract the token ID
	tokenID, err := parsedToken.GetString("token_id")
	if err != nil {
		return nil, err
	}

	uuidTokenID, err := uuid.Parse(tokenID)
	if err != nil {
		return nil, err
	}

	// Extract the subject (username)
	username, err := parsedToken.GetSubject()
	if err != nil {
		return nil, err
	}

	// Extract the issued at and expiration times
	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, err
	}

	expiredAt, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, err
	}

	// Return the payload
	return &Payload{
		ID:        uuidTokenID,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
