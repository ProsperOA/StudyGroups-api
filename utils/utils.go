package utils

import (
  "errors"
  "strings"
)

func GetAuthTokenFromHeader(header string) (string, error) {
  authHeader := strings.Split(header, " ")

  if len(authHeader) < 2 {
    return header, errors.New("invalid authorization header")
  }

  return authHeader[1], nil
}
