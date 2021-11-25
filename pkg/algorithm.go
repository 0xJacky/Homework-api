package pkg

import (
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return cast.ToString(bytes)
}

func PasswordVerify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
