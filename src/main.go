package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/handlers"
  "github.com/prosperoa/study-groups/src/server"
)

func main() {
  if err := server.InitServer(); err != nil {
    log.Fatal(err)
  }

  router := gin.Default()
  router.NoRoute(noRouteFound)

  public := router.Group("/api/v1")
  public.GET("/", index)
  public.POST("/login",  handlers.Login)
  public.POST("/signup", handlers.Signup)

  private := router.Group("/api/v1")
  private.Use(handlers.Auth())
  private.GET(   "/users",     handlers.GetUsers)
  private.GET(   "/users/:id", handlers.GetUser)
  private.DELETE("/users/:id", handlers.DeleteUser)

  log.Fatal(router.Run(":8080"))
}

func index(c *gin.Context) {
  server.Respond(c, nil, "StudyGroups API v1", http.StatusOK)
}

func noRouteFound(c *gin.Context) {
  server.Respond(c, nil, "route not found", http.StatusNotFound)
}
