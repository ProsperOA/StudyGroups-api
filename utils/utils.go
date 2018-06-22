package utils

import (
  "errors"
  "strings"
  "unicode"
)

func GetAuthTokenFromHeader(header string) (string, error) {
  authHeader := strings.Split(header, " ")

  if len(authHeader) < 2 {
    return header, errors.New("invalid authorization header")
  }

  return authHeader[1], nil
}

func IsInt(s string) bool {
  for _, c := range s {
    if !unicode.IsDigit(c) {
      return false
    }
  }

  return true
}
