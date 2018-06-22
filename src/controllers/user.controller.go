package controllers

import (
  "database/sql"
  "errors"
  "net/http"

  "github.com/prosperoa/study-groups/src/models"
  "github.com/prosperoa/study-groups/src/server"
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

func GetUsers(page, pageSize int) ([]models.User, int, error) {
  var users []models.User

  err := server.DB.Select(&users, "SELECT * FROM users LIMIT $1 OFFSET $2",
    pageSize, pageSize * page,
  )

  switch {
    case err == sql.ErrNoRows, len(users) == 0:
      return users, http.StatusNotFound, errors.New("no users found")
    case err != nil:
      return users, http.StatusInternalServerError, errors.New("unable to get users")
  }

  return users, http.StatusOK, nil
}
