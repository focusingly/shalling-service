package encrypt

import "golang.org/x/crypto/bcrypt"

func EncryptPasswordByBcrypt(rawPassword string) (string, error) {
	bf, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bf), err
}

func ComparePassword(rawPassword, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword)); err != nil {
		return false
	} else {
		return true
	}
}
