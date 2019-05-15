package auth

import (
	"fmt"
	"strings"
)

const (
	AuthHeader     = "Authorization"
	BearerAuthType = "Bearer"
	Separator      = " "
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

func ParseAuth(authHeader string) error {
	authHeaderParts := strings.Split(authHeader, Separator)
	if len(authHeaderParts) <= 1 {
		return NewBadRequestError("invalid auth header format. expected '<type> <credentials>'")
	}

	authType := authHeaderParts[0]
	authString := strings.Join(authHeaderParts[1:], Separator)

	switch authType {
	case BearerAuthType:
		return DecodeToken(authString)
	}

	return NewBadRequestError(
		fmt.Sprintf("invalid auth type: `%v`", authType),
	)
}
