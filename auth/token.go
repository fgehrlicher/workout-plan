package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type UserAccessClaim struct {
	Grant
	jwt.StandardClaims
}

func DecodeToken(rawToken string) error {

	return nil
}
