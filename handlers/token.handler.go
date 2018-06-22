package handlers

import (
  "errors"
  "time"

  "github.com/dgrijalva/jwt-go"
)

var (
  signingKey = []byte("secret")
)

func GenerateAuthToken() (string, error) {
  token := jwt.New(jwt.SigningMethodHS256)

  token.Claims = jwt.MapClaims{
    "exp": time.Now().Add(time.Hour * 730).Unix(), // ~ 1 month
    "iat": time.Now().Unix(),
  }

  tokenString, err := token.SignedString(signingKey)

  if err != nil { return tokenString, errors.New("error while signing auth token") }

  return tokenString, nil
}

func VerifyAuthToken(t string) error {
  errMsg := errors.New("invalid auth token")

  authToken, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errMsg
    }

    if int64(token.Claims.(jwt.MapClaims)["exp"].(float64)) < time.Now().Unix() {
      return nil, errMsg
    }

    return signingKey, nil
  })

  if err == nil && authToken.Valid {
    return nil
  } else {
    return errMsg
  }
}
