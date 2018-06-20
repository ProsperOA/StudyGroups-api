package server

import (
  "github.com/gin-gonic/gin"
)

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
