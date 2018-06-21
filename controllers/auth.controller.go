package controllers

import (
  "database/sql"
  "errors"
  "net/http"

  "github.com/prosperoa/study-groups/models"
  "github.com/prosperoa/study-groups/server"
  "golang.org/x/crypto/bcrypt"
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

  if email != user.Email {
    return user, http.StatusBadRequest, errors.New("incorrect email")
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

  return user, http.StatusOK, nil
}
