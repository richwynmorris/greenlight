package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"math/rand"
	"time"

	"richwynmorris.co.uk/internal/validator"
)

const (
	scopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the token in plaintext to be sent back to the user to authenticate with.
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generates a SHA256 of the plain text token to be stored in the tokens db. When the user trys to authenticate
	// we'll check the received token against the hashed token and validate if the user is authentic.
	hash := sha256.Sum256([]byte(token.Plaintext))
	// Convert returned array into a slice for easier use.
	token.Hash = hash[:]

	return token, nil
}

// =========================== TOKEN VALIDATION ==========================================

func validateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// ========================= TOKEN DATABASE MODEL =======================================

type TokenModel struct {
	DB *sql.DB
}

// New generates a new token and inserts it into the tokens database.
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
			  VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`

	args := []any{scope, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
