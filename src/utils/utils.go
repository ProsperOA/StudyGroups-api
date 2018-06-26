package utils

import (
  "errors"
  "math/rand"
  "strings"
  "time"
  "unicode"
)

const (
  KB int64 = 1 << (10 * iota)
  MB
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
  rand.Seed(time.Now().UnixNano())
}

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

func RandString(n int) string {
  b := make([]rune, n)

  for i := range b {
    b[i] = letterRunes[rand.Intn(len(letterRunes))]
  }

  return string(b)
}
