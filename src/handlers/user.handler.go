package handlers

import (
  "net/http"
  "path/filepath"
  "strconv"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/controllers"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/utils"
)

func GetUser(c *gin.Context) {
  userID := c.Param("id")

  if userID == "" || !utils.IsInt(userID) {
    server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
    return
  }

  user, status, err := controllers.GetUser(userID)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }


  server.Respond(c, user, "", status)
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

  server.Respond(c, users, "", status)
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

func GetUserStudyGroups(c *gin.Context) {
  userID   := c.Param("id")
  page     := c.DefaultQuery("page", "0")
  pageSize := c.DefaultQuery("page_size", "30")

  if userID == "" || page == "" || pageSize == "" ||
    !utils.IsInt(userID) || !utils.IsInt(page) || !utils.IsInt(pageSize) {
      server.Respond(c, nil, "invalid params", http.StatusBadRequest)
      return
  }

  p, _  := strconv.Atoi(page)
  ps, _ := strconv.Atoi(pageSize)

  studyGroups, status, err := controllers.GetUserStudyGroups(userID, p, ps)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, studyGroups, "", status)
}

func UploadAvatar(c *gin.Context) {
  userID := c.Param("id")
  file, err := c.FormFile("image")

  if userID == "" ||  !utils.IsInt(userID) {
    server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
    return
  }

  if err != nil {
    server.Respond(c, nil, "unable to get image", http.StatusBadRequest)
    return
  }

  if file.Size / utils.MB > utils.MB * 2 {
    server.Respond(c, nil, "image size must be 2MB or less", http.StatusBadRequest)
    return
  }

  ext := filepath.Ext(file.Filename)
  image, _ := file.Open()
  avatarURL, status, err := controllers.UploadAvatar(userID, ext, image)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
  }

  server.Respond(c, avatarURL, "", status)
}
