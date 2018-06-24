package controllers

import (
  "database/sql"
  "errors"
  "log"
  "net/http"

  "github.com/prosperoa/study-groups/src/models"
  "github.com/prosperoa/study-groups/src/server"
)

func GetStudyGroup(id string) (models.StudyGroup, int, error) {
  var studyGroup models.StudyGroup

  err := server.DB.Get(&studyGroup, "SELECT * FROM study_groups WHERE id = $1", id)

  switch {
    case err == sql.ErrNoRows:
      return studyGroup, http.StatusNotFound, errors.New("study group not found")
    case err != nil:
      return studyGroup, http.StatusInternalServerError, errors.New(
        "unable to get study group",
      )
  }

  return studyGroup, http.StatusOK, nil
}

func GetStudyGroups(page, pageSize int) ([]models.StudyGroup, int, error) {
  var studyGroups []models.StudyGroup

  err := server.DB.Select(&studyGroups, "SELECT * FROM study_groups LIMIT $1 OFFSET $2",
    pageSize, pageSize * page,
  )

  switch {
    case err == sql.ErrNoRows, len(studyGroups) == 0:
      return studyGroups, http.StatusNotFound, errors.New("no study groups found")
    case err != nil:
      return studyGroups, http.StatusInternalServerError, errors.New(
        "unable to get study groups",
      )
  }

  return studyGroups, http.StatusOK, nil
}

func DeleteStudyGroup(id string) (int, error) {
  result, err := server.DB.Exec("DELETE FROM study_groups WHERE id = $1", id)
  rowsAffected, _ := result.RowsAffected()

  switch {
    case rowsAffected == 0:
      return http.StatusBadRequest, errors.New("study group doesn't exist")
    case err != nil:
      return http.StatusInternalServerError, errors.New("unable delete study group")
  }

  return http.StatusOK, nil
}

func LeaveStudyGroup(studyGroupID, userID string) (int, error) {
  var user models.User
  var studyGroup models.StudyGroup
  var userStudyGroups string
  var studyGroupMembers string

  internalErr := func() (int, error) {
    return http.StatusInternalServerError, errors.New("unable to leave study group")
  }

  err := server.DB.Get(&studyGroup, "SELECT * FROM study_groups WHERE id = $1", studyGroupID)
  switch {
    case err == sql.ErrNoRows:
      return http.StatusNotFound, errors.New("study group not found")
    case err != nil:
      return internalErr()
  }

  err = server.DB.Get(&user.StudyGroups, "SELECT study_groups FROM users WHERE id = $1", userID)
  switch {
    case err == sql.ErrNoRows:
      return http.StatusNotFound, errors.New("user not found")
    case err != nil:
      return internalErr()
  }

  if err = studyGroup.RemoveMember(userID); err != nil { return internalErr() }
  if err = user.LeaveStudyGroup(studyGroupID); err != nil { return internalErr() }

  {
    tx, err := server.DB.Begin()
    if err != nil { return internalErr() }

    defer func() (int, error) {
      if err != nil {
        log.Println(err.Error())
        tx.Rollback()
        return internalErr()
      }

      return 0, nil
    }()

    studyGroupMembers = studyGroup.Members.String

    _, err = tx.Exec(
      `UPDATE study_groups
       SET members = $1, available_spots = available_spots + 1
       WHERE id = $2`,
       studyGroupMembers,
       studyGroupID,
     )

    userStudyGroups = user.StudyGroups.String

    _, err = tx.Exec("UPDATE users SET study_groups = $1 WHERE id = $2",
      userStudyGroups,
      userID,
    )

    err = tx.Commit()
  }

  return http.StatusOK, nil
}
