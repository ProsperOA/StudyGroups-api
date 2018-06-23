package middlewares

import (
  "errors"
  "net/http"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/utils"
)

func AuthHandler() gin.HandlerFunc {
  return func(c *gin.Context) {
    authToken, err := utils.GetAuthTokenFromHeader(c.GetHeader("Authorization"))

    if err != nil {
      server.Respond(c, nil, err.Error(), http.StatusBadRequest)
      c.Abort()
      return
    }

    err = verifyAuthToken(authToken)

    if err != nil {
      server.Respond(c, nil, err.Error(), http.StatusUnauthorized)
      c.Abort()
      return
    }

    c.Next()
  }
}

func verifyAuthToken(t string) error {
  errMsg := errors.New("invalid auth token")

  authToken, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errMsg
    }

    if int64(token.Claims.(jwt.MapClaims)["exp"].(float64)) < time.Now().Unix() {
      return nil, errMsg
    }

    return server.JWTSigningKey, nil
  })

  if err == nil && authToken.Valid {
    return nil
  } else {
    return errMsg
  }
}
