package server

import (
  "database/sql"
  "os"

  "github.com/gin-gonic/gin"
  _ "github.com/lib/pq"
)

type Context struct {
  DB *sql.DB
}

func (c *Context) InitServer() error {
  var err error

  c.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URI"))

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
