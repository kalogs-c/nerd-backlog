package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	ArgonTime    = 1
	ArgonMemory  = 64 * 1024
	ArgonThreads = 4
	ArgonKeyLen  = 32
	SaltLen      = 16
)

func genSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var res byte
	for i := range len(a) {
		res |= a[i] ^ b[i]
	}

	return res == 0
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password required")
	}
	salt, err := genSalt(SaltLen)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, ArgonTime, ArgonMemory, ArgonThreads, ArgonKeyLen)
	return fmt.Sprintf("%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

func ComparePassword(password string, hashedPassword string) (bool, error) {
	parts := strings.SplitN(hashedPassword, "$", 2)
	if len(parts) != 2 {
		return false, ErrInvalidHashFormat
	}

	saltB64 := parts[0]
	hashB64 := parts[1]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), salt, ArgonTime, ArgonMemory, ArgonThreads, uint32(len(expectedHash)))
	if subtleCompare(hash, expectedHash) {
		return true, nil
	}

	return false, nil
}
