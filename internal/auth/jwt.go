package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/karrick/tparse"
	"github.com/pkg/errors"
	"incrementor/internal/config"
	"incrementor/internal/models"
	"time"
)

// Key type fot context key used to pass username
type Key string

const (
	// ContextKey const for context key
	ContextKey = Key("Username")
)

// Jwt hold sturct for token generateion and validation
type Jwt struct {
	Secret  string
	Expired time.Time
	Claims  *Claims
}

// Claims for JWT
type Claims struct {
	Username string
	jwt.StandardClaims
}

// NewJwt construct Jwt instance
func NewJwt(config *config.Config) (*Jwt, error) {

	now := time.Now()
	expired, err := tparse.AddDuration(now, config.JWT.Duration)
	if err != nil {
		return nil, errors.Wrap(err, "Error parse duration for JWT token")
	}

	j := &Jwt{
		Secret:  config.JWT.Secret,
		Expired: expired,
	}

	return j, nil
}

// Sign creates new jwt token for Client
func (j *Jwt) Sign(client *models.Client) (string, error) {

	j.Claims = &Claims{
		Username: client.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: j.Expired.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *j.Claims)

	tokenString, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", errors.Wrap(err, "Error signing JWT token")
	}

	return tokenString, nil
}

// Validate validate jwt token from string
func (j *Jwt) Validate(tokenStr string) (bool, error) {

	j.Claims = &Claims{}

	if tokenStr == "" {
		return false, errors.New("Empty token")
	}

	tkn, err := jwt.ParseWithClaims(tokenStr, j.Claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if false == tkn.Valid {
		return false, errors.New("Invalid token")
	}

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, errors.New("Jwt invalid signature")
		}
		return false, errors.New("Bad jwt token")
	}

	return true, nil
}
