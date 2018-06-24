package handlers

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/controllers"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/utils"
)

func GetStudyGroup(c *gin.Context) {
  id := c.Param("id")

  if id == "" || !utils.IsInt(id) {
    server.Respond(c, nil, "invalid study group id", http.StatusBadRequest)
    return
  }

  studyGroup, status, err := controllers.GetStudyGroup(id)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }


  server.Respond(c, studyGroup, "", status)
}

func GetStudyGroups(c *gin.Context) {
  page     := c.DefaultQuery("page", "0")
  pageSize := c.DefaultQuery("page_size", "30")

  if page == "" || pageSize == "" ||
    !utils.IsInt(page) || !utils.IsInt(pageSize) {
      server.Respond(c, nil, "invalid params", http.StatusBadRequest)
      return
  }

  p, _  := strconv.Atoi(page)
  ps, _ := strconv.Atoi(pageSize)

  studyGroups, status, err := controllers.GetStudyGroups(p, ps)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, studyGroups, "", status)
}

func JoinStudyGroup(c *gin.Context) {
  studyGroupID := c.Param("id")
  userID := c.PostForm("user_id")

  if studyGroupID == "" || userID == "" ||
    !utils.IsInt(studyGroupID) || !utils.IsInt(userID) {
      server.Respond(c, nil, "invalid params", http.StatusBadRequest)
      return
  }

  status, err := controllers.JoinStudyGroup(studyGroupID, userID)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, nil, "user added to study group waitlist", status)
}

func DeleteStudyGroup(c *gin.Context) {
  id := c.Param("id")

  if id == "" || !utils.IsInt(id) {
    server.Respond(c, nil, "invalid study group id", http.StatusBadRequest)
    return
  }

  status, err := controllers.DeleteStudyGroup(id)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, nil, "study group successfully deleted", status)
}

func LeaveStudyGroup(c *gin.Context) {
  studyGroupID := c.Param("id")
  userID := c.PostForm("user_id")

  if studyGroupID == "" || userID == "" ||
    !utils.IsInt(studyGroupID) || !utils.IsInt(userID) {
      server.Respond(c, nil, "invalid study group id", http.StatusBadRequest)
      return
  }

  status, err := controllers.LeaveStudyGroup(studyGroupID, userID)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  server.Respond(c, nil, "user removed from study group", status)
}
