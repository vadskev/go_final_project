package env

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
)

const (
	PasswordEnvName = "TODO_PASSWORD"
)

type Password interface {
	GetPass() string
	CreateHash(key string) string
}

type passConfig struct {
	password string
}

func NewPassConfig() (Password, error) {
	pass := os.Getenv(PasswordEnvName)
	if len(pass) == 0 {
		return nil, errors.New("password not found")
	}

	return &passConfig{
		password: pass,
	}, nil
}

func (cfg *passConfig) GetPass() string {
	return cfg.password
}

func (cfg *passConfig) CreateHash(key string) string {
	h := hmac.New(sha256.New, []byte(cfg.GetPass()))
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
