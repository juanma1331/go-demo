package infra

import (
	"golang.org/x/crypto/bcrypt"
)

type bCryptPasswordManager struct{}

func NewBCryptPasswordManager() *bCryptPasswordManager {
	return &bCryptPasswordManager{}
}

func (pm bCryptPasswordManager) GenerateFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (pm bCryptPasswordManager) CompareHashAndPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
