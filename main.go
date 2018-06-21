package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  // "github.com/prosperoa/study-groups/handlers"
  "github.com/prosperoa/study-groups/server"
)

func main() {
  var s server.Context

  if err := s.InitServer(); err != nil {
    panic(err)
  }

  router := gin.Default()

  router.GET("/", index)

  log.Fatal(router.Run(":8080"))
}

func index(c *gin.Context) {
  server.Respond(c, nil, "StudyGroups API", http.StatusOK)
}
