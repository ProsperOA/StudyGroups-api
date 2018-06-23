package server

import (
  "errors"
  "os"

  "github.com/badoux/checkmail"
  "github.com/gin-gonic/gin"
  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
)

var (
  DB *sqlx.DB
)

func InitServer() error {
  var err error

  DB, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URI"))

  if err != nil { return err }

  return nil
}

func ValidateEmail(email string) (errMsg error) {
  errMsg = errors.New("invalid email address")

  if err := checkmail.ValidateFormat(email); err != nil {
    return errMsg
  }

  if err := checkmail.ValidateHost(email); err != nil {
    return errMsg
  }

  err := checkmail.ValidateHost(email)
  if _, ok := err.(checkmail.SmtpError); ok && err != nil {
    return errMsg
  }

  return nil
}

func Respond(c *gin.Context, data interface{}, message string, httpStatus int) {
  var success bool

  // HTTP Status Code >= 400 indicates error
  if httpStatus < 400 {
    success = true
  } else {
    success = false
  }

  c.JSON(
    httpStatus,
    map[string]interface{} {
      "data":    data,
      "status":  httpStatus,
      "message": message,
      "success": success,
    },
  )
}
