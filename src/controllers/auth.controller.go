package controllers

import (
  "database/sql"
  "errors"
  "log"
  "net/http"
  "os"

  mailchimp "github.com/beeker1121/mailchimp-go"
  "github.com/beeker1121/mailchimp-go/lists/members"
  "github.com/prosperoa/study-groups/src/models"
  "github.com/prosperoa/study-groups/src/server"
  "golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (models.User, int, error) {
  var user models.User

  err := server.DB.Get(&user, "SELECT * FROM users WHERE email = $1", email)

  switch {
    case err == sql.ErrNoRows:
      return user, http.StatusUnauthorized, errors.New("account does not exist")
    case err != nil:
      return user, http.StatusInternalServerError, errors.New("unable to login")
  }

  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
  if err != nil {
    return user, http.StatusBadRequest, errors.New("incorrect password")
  }

  return user, http.StatusOK, nil
}

func Signup(firstName, lastName, email, password string) (models.User, int, error) {
  var user models.User
  var accountExists bool
  errMsg := "unable to create account"

  err := server.DB.Get(&accountExists, "SELECT exists(SELECT 1 FROM users WHERE email = $1)",
    email,
  )

  if accountExists {
    return user, http.StatusForbidden, errors.New("account already exists")
  } else if err != nil {
    return user, http.StatusInternalServerError, errors.New(errMsg)
  }

  passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 6)

  if err != nil {
    return user, http.StatusInternalServerError, errors.New(errMsg)
  }

  err = server.DB.Get(
    &user,
    `INSERT INTO users (first_name, last_name, email, password)
     VALUES ($1, $2, $3, $4) RETURNING *`,
     firstName, lastName, email, passwordHash,
  )

  if err != nil {
    return user, http.StatusInternalServerError, errors.New(errMsg)
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

  return user, http.StatusOK, nil
}
