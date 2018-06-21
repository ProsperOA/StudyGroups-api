package controllers

import (
  "database/sql"
  "errors"
  "net/http"

  "github.com/prosperoa/study-groups/models"
  "github.com/prosperoa/study-groups/server"
)

func Login(userID, email, password string) (models.User, int, error) {
  var user models.User

  err := server.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userID)

  switch {
    case err == sql.ErrNoRows:
      return user, http.StatusUnauthorized, errors.New("account does not exist")
    case err != nil:
      return user, http.StatusInternalServerError, errors.New("unable to login")
  }

  if email != user.Email || password != user.Password {
    return user, http.StatusUnauthorized, errors.New("invalid email or password")
  }

  return user, http.StatusOK, nil
}
