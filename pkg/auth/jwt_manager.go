package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTManager struct {
	Secret     []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type claims struct {
	AccountID uuid.UUID `json:"account_id"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret []byte, accessTTL, refreshTTL time.Duration) JWTManager {
	return JWTManager{
		Secret:     secret,
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}
}

func (j *JWTManager) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now().UTC()

	c := claims{
		AccountID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.AccessTTL)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(j.Secret)
}

func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(j.RefreshTTL)

	c := claims{
		AccountID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := t.SignedString(j.Secret)
	return signed, exp, err
}

func (j *JWTManager) VerifyAccessToken(token string) (uuid.UUID, error) {
	parsed, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.Secret, nil
	})
	if err != nil || !parsed.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	parsedClaims, ok := parsed.Claims.(*claims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return parsedClaims.AccountID, nil
}
