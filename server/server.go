package server

import (
  "os"

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
