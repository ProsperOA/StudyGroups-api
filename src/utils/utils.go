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

func Splice(s []string, elem string) []string {
  for i, v := range s {
    if v == elem {
      s = append(s[:i], s[i + 1:]...)
      break
    }
  }

  return s
}

func Contains(s []string, elem string) bool {
  for _, v := range s {
    if v == elem { return true }
  }

  return false
}
