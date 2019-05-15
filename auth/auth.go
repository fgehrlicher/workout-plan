package auth

import (
	"errors"
	"fmt"
	"strings"
)

const (
	AuthorizationHeader     = "Authorization"
	BearerAuthorizationType = "Bearer"
	Separator               = " "
)

func ParseAuth(authorizationHeader string) error {
	authorizationHeaderParts := strings.Split(authorizationHeader, Separator)
	if len(authorizationHeaderParts) <= 1 {
		return errors.New("invalid authorization header format. expected '<type> <credentials>'")
	}

	authorizationType := authorizationHeaderParts[0]
	authorizationString := strings.Join(authorizationHeaderParts[1:], Separator)

	switch authorizationType {
	case BearerAuthorizationType:
		return DecodeToken(authorizationString)
	}

	return fmt.Errorf("invalid authorization type: `%v`", authorizationType)
}
