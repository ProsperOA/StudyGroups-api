package controllers

import (
  "database/sql"
  "errors"
  "net/http"

  "github.com/prosperoa/study-groups/models"
  "github.com/prosperoa/study-groups/server"
)

func GetUser(userID string) (models.User, int, error) {
  var user models.User

  err := server.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userID)

  switch {
    case err == sql.ErrNoRows:
      return user, http.StatusNotFound, errors.New("user not found")
    case err != nil:
      return user, http.StatusInternalServerError, errors.New("unable to get user")
  }

  return user, http.StatusOK, nil
}
