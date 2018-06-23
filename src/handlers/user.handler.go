package handlers

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/controllers"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/utils"
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

func GetUsers(c *gin.Context) {
  page     := c.DefaultQuery("page", "0")
  pageSize := c.DefaultQuery("page_size", "30")

  if page == "" || pageSize == "" ||
    !utils.IsInt(page) || !utils.IsInt(pageSize) {
      server.Respond(c, nil, "invalid params", http.StatusBadRequest)
      return
  }

  p, _  := strconv.Atoi(page)
  ps, _ := strconv.Atoi(pageSize)

  users, status, err := controllers.GetUsers(p, ps)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, users, "", http.StatusOK)
}

func DeleteUser(c *gin.Context) {
  userID := c.Param("id")

  if userID == "" || !utils.IsInt(userID) {
    server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
    return
  }

  status, err := controllers.DeleteUser(userID)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, nil, "account successfully deleted", status)
}
