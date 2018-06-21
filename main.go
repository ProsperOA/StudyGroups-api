package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/handlers"
  "github.com/prosperoa/study-groups/server"
)

func main() {
  if err := server.InitServer(); err != nil {
    log.Fatal(err)
  }

  router := gin.Default()

  router.NoRoute(noRouteFound)

  router.GET("/", index)

  router.POST("/login",  handlers.Login)
  router.POST("/signup", handlers.Signup)

  log.Fatal(router.Run(":8080"))
}

func index(c *gin.Context) {
  server.Respond(c, nil, "StudyGroups API", http.StatusOK)
}

func noRouteFound(c *gin.Context) {
  server.Respond(c, nil, "page not found", http.StatusNotFound)
}
