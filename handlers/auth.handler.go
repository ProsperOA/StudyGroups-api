package handlers

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/controllers"
  "github.com/prosperoa/study-groups/server"
)

func Login(c *gin.Context) {
  userID   := c.PostForm("user_id")
  email    := c.PostForm("email")
  password := c.PostForm("password")

  if userID == "" || email == "" || password == "" {
    server.Respond(c, nil, "invalid params", http.StatusBadRequest)
    return
  }

  user, status, err := controllers.Login(userID, email, password)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, user, "", status)
}
