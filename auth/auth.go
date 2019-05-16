package auth

import (
	"fmt"
	"strings"
)

const (
	Header     = "Authorization"
	BearerType = "Bearer"
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
	authParts := strings.Split(authHeader, " ")
	if len(authParts) <= 1 {
		return NewBadRequestError("invalid auth header format. expected '<type> <credentials>'")
	}

	authType := authParts[0]
	authString := strings.Join(authParts[1:], " ")

	switch authType {
	case BearerType:
		return DecodeToken(authString)
	}

	return NewBadRequestError(
		fmt.Sprintf("invalid auth type: `%v`", authType),
	)
}
