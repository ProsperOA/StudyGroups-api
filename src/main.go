package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/handlers"
  "github.com/prosperoa/study-groups/src/middlewares"
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
  private.Use(middlewares.AuthHandler())
  private.GET(   "/users",     handlers.GetUsers)
  private.GET(   "/users/:id", handlers.GetUser)
  private.DELETE("/users/:id", handlers.DeleteUser)

  private.GET(   "/study_groups",           handlers.GetStudyGroups)
  private.GET(   "/study_groups/:id",       handlers.GetStudyGroup)
  private.DELETE("/study_groups/:id",       handlers.DeleteStudyGroup)
  private.PATCH( "/study_groups/:id/leave", handlers.LeaveStudyGroup)

  log.Fatal(router.Run(":8080"))
}

func index(c *gin.Context) {
  server.Respond(c, nil, "StudyGroups API v1", http.StatusOK)
}

func noRouteFound(c *gin.Context) {
  server.Respond(c, nil, "route not found", http.StatusNotFound)
}
