package middlewares

import (
	"errors"
	"net/http"
	"strings"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/utils"
)

func BasicAuth() gin.HandlerFunc {
  return func(c *gin.Context) {
		authToken, err := utils.GetAuthTokenFromHeader(c.GetHeader("Authorization"))

    if err != nil {
      server.Respond(c, nil, err.Error(), http.StatusBadRequest)
      c.Abort()
      return
    }

    if err = verifyBasicAuth(authToken); err != nil {
      server.Respond(c, nil, err.Error(), http.StatusUnauthorized)
      c.Abort()
      return
    }

    c.Next()
  }
}

func ResourceOwnerAuth() gin.HandlerFunc {
  return func(c *gin.Context) {
		userID := c.Param("id")

		if userID == "" ||
			!strings.Contains(c.Request.URL.String(), "users") ||
			c.Request.Method == "GET" {
				c.Next()
				return
		}

    authToken, err := utils.GetAuthTokenFromHeader(c.GetHeader("Authorization"))

    if err != nil {
      server.Respond(c, nil, err.Error(), http.StatusBadRequest)
      c.Abort()
      return
    }

    if err = verifyResourceOwnerAuth(authToken, userID); err != nil {
      server.Respond(c, nil, err.Error(), http.StatusUnauthorized)
      c.Abort()
      return
    }

    c.Next()
  }
}

func verifyBasicAuth(t string) error {
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

  if err != nil || !authToken.Valid {	return err }

	return nil
}

func verifyResourceOwnerAuth(t, userID string) error {
  errMsg := errors.New("resource access unauthorized")

  authToken, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errMsg
		}

		claims := token.Claims.(jwt.MapClaims)
    if claims["user_id"] != userID {
      return nil, errMsg
    }

    return server.JWTSigningKey, nil
  })

  if err != nil || !authToken.Valid { return errMsg }

	return nil
}
