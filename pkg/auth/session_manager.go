package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

const SessionCookieName = "nb_session"

const sessionTokenBytes = 32

type SessionManager struct {
	TTL time.Duration
}

func NewSessionManager(ttl time.Duration) SessionManager {
	return SessionManager{TTL: ttl}
}

func (s *SessionManager) GenerateSessionToken() (string, time.Time, error) {
	buf := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", time.Time{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(buf)
	expiresAt := time.Now().UTC().Add(s.TTL)
	return token, expiresAt, nil
}
