package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"workout-plan/config"
)

const algHeader = "alg"

type UserAccessClaim struct {
	Grant
	jwt.StandardClaims
}

func DecodeToken(rawToken string, config config.TokenConfig) (*Grant, error) {
	hmacSecret := []byte(config.Secret)

	parsedToken, err := jwt.ParseWithClaims(rawToken, &UserAccessClaim{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if ok {
			return hmacSecret, nil
		}

		return nil, NewBadRequestError(
			fmt.Sprintf("unexpected signing method: %v", token.Header[algHeader]),
		)
	})

	if err != nil {
		return nil, err
	}

	claim, ok := parsedToken.Claims.(*UserAccessClaim)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %v", parsedToken)
	}

	claim.UserName = claim.Subject
	err = claim.Validate(config)
	return &claim.Grant, err
}

func (claim *UserAccessClaim) Validate(config config.TokenConfig) error {
	err := claim.Valid()
	if err != nil {
		return err
	}

	if !claim.VerifyAudience(config.Service, true) {
		return fmt.Errorf("invalid token audience")
	}

	if !claim.VerifyIssuer(config.Issuer, true) {
		return fmt.Errorf("invalid token issuer")
	}

	return nil
}
