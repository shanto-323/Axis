package pkg

import "golang.org/x/crypto/bcrypt"

func CreateHash(passwordString string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordString), 10)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CompareWithHash(passwordHash, passwordString string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordString))
}

