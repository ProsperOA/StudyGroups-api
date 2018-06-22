package handlers

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/controllers"
  "github.com/prosperoa/study-groups/server"
)

func GetUser(c *gin.Context) {
  userID := c.Param("id")

  if userID == "" {
    server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
    return
  }

  user, status, err := controllers.GetUser(userID)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }


  server.Respond(c, user, "", http.StatusOK)
}
