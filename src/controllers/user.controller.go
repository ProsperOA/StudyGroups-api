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

func DeleteUser(userID string) (int, error) {
  var userEmail string

  err := server.DB.Get(&userEmail, "DELETE FROM users WHERE id = $1 RETURNING email", userID)

  switch {
    case err == sql.ErrNoRows:
      return http.StatusBadRequest, errors.New("account doesn't exist")
    case err != nil:
      return http.StatusInternalServerError, errors.New("unable delete account")
  }

  // remove user from mailchimp list
  var mailchimpID string

  params := &members.GetParams{Status: members.StatusSubscribed}
  listMembers, err := members.Get("4d6392ba4d", params)

  if err != nil {
    for _, v := range listMembers.Members {
      if v.EmailAddress == userEmail {
        mailchimpID = v.ID
        break
      }
    }

    if err = mailchimp.SetKey(os.Getenv("MAILCHIMP_API_KEY")); err != nil {
      log.Println(err.Error())
    }

    if err = members.Delete("4d6392ba4d", mailchimpID); err != nil {
      log.Println(err.Error())
    }
  } else {
    log.Println(err.Error())
  }

  return http.StatusOK, nil
}

func GetUserStudyGroups(userID string, page, pageSize int) ([]models.StudyGroup, int, error) {
  var studyGroups []models.StudyGroup

  err := server.DB.Select(
    &studyGroups,
    "SELECT * FROM study_groups WHERE user_id = $1 LIMIT $2 OFFSET $3",
    userID,
    pageSize,
    pageSize * page,
  )

  switch {
    case err == sql.ErrNoRows, len(studyGroups) == 0:
      return studyGroups, http.StatusNotFound, errors.New("no users study groups found")
    case err != nil:
      return studyGroups, http.StatusInternalServerError, errors.New("unable to get user's study groups")
  }

  return studyGroups, http.StatusOK, nil
}
