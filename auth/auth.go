package auth

import (
	"fmt"
	"strings"
)

const (
	AuthorizationHeader     = "Authorization"
	BearerAuthorizationType = "Bearer"
	Separator               = " "
)

func NewBadRequestError(text string) error {
	return &BadRequestError{text}
}

type BadRequestError struct {
	text string
}

func (badRequestError *BadRequestError) Error() string {
	return badRequestError.text
}

func ParseAuth(authorizationHeader string) error {
	authorizationHeaderParts := strings.Split(authorizationHeader, Separator)
	if len(authorizationHeaderParts) <= 1 {
		return NewBadRequestError("invalid authorization header format. expected '<type> <credentials>'")
	}

	authorizationType := authorizationHeaderParts[0]
	authorizationString := strings.Join(authorizationHeaderParts[1:], Separator)

	switch authorizationType {
	case BearerAuthorizationType:
		return DecodeToken(authorizationString)
	}

	return NewBadRequestError(
		fmt.Sprintf("invalid authorization type: `%v`", authorizationType),
	)
}
