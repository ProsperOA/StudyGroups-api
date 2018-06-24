package handlers

import (
  "log"
  "net/http"
  "os"

  mailchimp "github.com/beeker1121/mailchimp-go"
  "github.com/beeker1121/mailchimp-go/lists/members"
  "github.com/gin-gonic/gin"
  "github.com/prosperoa/study-groups/src/controllers"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/email-notifications"
)

func Login(c *gin.Context) {
  email    := c.PostForm("email")
  password := c.PostForm("password")

  if email == "" || password == "" {
    server.Respond(c, nil, "invalid params", http.StatusBadRequest)
    return
  }

  user, status, err := controllers.Login(email, password)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  authToken, err := server.GenerateAuthToken()
  if err != nil {
    server.Respond(c, nil, err.Error(), http.StatusInternalServerError)
    return
  }

  data := map[string]interface{} {
    "auth_token": authToken,
    "user": user,
  }

  server.Respond(c, data, "", status)
}

func Signup(c *gin.Context) {
  firstName := c.Query("first_name")
  lastName  := c.Query("last_name")
  email     := c.Query("email")
  password  := c.Query("password")

  if firstName == "" ||  email == "" || password == "" {
    server.Respond(c, nil, "invalid params", http.StatusBadRequest)
    return
  }

  if err := server.ValidateEmail(email); err != nil {
    server.Respond(c, nil, err.Error(), http.StatusBadRequest)
    return
  }

  if len(password) < 6 || len(password) > 50 {
    server.Respond(c, nil, "password must contain 6 to 50 characters",
      http.StatusBadRequest,
    )
    return
  }

  authToken, err := server.GenerateAuthToken()
  if err != nil {
    server.Respond(c, nil, err.Error(), http.StatusInternalServerError)
    return
  }

  user, status, err := controllers.Signup(firstName, lastName, email, password)

  if err != nil {
    server.Respond(c, nil, err.Error(), status)
    return
  }

  // add user to mailchimp list
  if err = mailchimp.SetKey(os.Getenv("MAILCHIMP_API_KEY")); err != nil {
    log.Println(err.Error())
  }

  params := &members.NewParams {
    EmailAddress: user.Email,
    Status: members.StatusSubscribed,
  }

  _, err = members.New("4d6392ba4d", params)
  if err != nil {
    log.Println(err.Error())
  }

  if err = emails.NewUserNotification(user.FirstName, user.Email); err != nil {
    log.Println(err.Error())
  }

  data := map[string]interface{} {
    "auth_token": authToken,
    "user": user,
  }

  server.Respond(c, data, "", status)
}
