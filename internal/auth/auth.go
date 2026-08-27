package auth

import "github.com/alexedwards/argon2id"

func GenerateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CompareHash(password string, passHash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, passHash)
}
