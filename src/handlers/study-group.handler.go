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
