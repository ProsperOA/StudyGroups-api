package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/controllers"
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
	private.Use(middlewares.BasicAuth())

  private.GET(   "/users",                  controllers.GetUsers)
  private.GET(   "/users/:id",              controllers.GetUser)
  private.PATCH( "/users/:id/account",      controllers.UpdateAccount)
  private.POST(  "/users/:id/avatar",       controllers.UploadAvatar)
  private.PUT(   "/users/:id/courses",      controllers.UpdateCourses)
  private.POST(  "/users/:id/delete",       controllers.DeleteUser)
  private.PATCH( "/users/:id/password",     controllers.ChangePassword)
  // private.GET(   "/users/:id/study_groups", handlers.GetUserStudyGroups)

  private.GET(   "/study_groups",                         handlers.GetStudyGroups)
  private.POST(  "/study_groups",                         handlers.CreateStudyGroup)
  private.GET(   "/study_groups/:id",                     handlers.GetStudyGroup)
  private.PATCH( "/study_groups/:id",                     handlers.UpdateStudyGroup)
  private.POST(  "/study_groups/:id",                     handlers.DeleteStudyGroup)
  private.POST(  "/study_groups/:id/join",                handlers.JoinStudyGroup)
  private.PATCH( "/study_groups/:id/leave",               handlers.LeaveStudyGroup)
  private.GET(   "/study_groups/:id/members",             handlers.GetStudyGroupMembers)
  private.PATCH( "/study_groups/:id/waitlist_to_members", handlers.MoveUserFromWaitlistToMembers)

  log.Fatal(router.Run(":8080"))
}

func index(c *gin.Context) {
  server.Respond(c, nil, "StudyGroups API v1", http.StatusOK)
}

func noRouteFound(c *gin.Context) {
  server.Respond(c, nil, "route not found", http.StatusNotFound)
}
