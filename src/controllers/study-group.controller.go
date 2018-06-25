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

func DeleteStudyGroup(id, userID string) (int, error) {
  result, err := server.DB.Exec("DELETE FROM study_groups WHERE id = $1 AND user_id = $2",
    id, userID,
  )
  rowsAffected, _ := result.RowsAffected()

  switch {
    case rowsAffected == 0:
      return http.StatusForbidden, errors.New("user can only delete their own study groups")
    case err != nil:
      return http.StatusInternalServerError, errors.New("unable delete study group")
  }

  return http.StatusOK, nil
}

func JoinStudyGroup(studyGroupID, userID string) (int, error) {
  var user models.User
  var studyGroup models.StudyGroup

  internalErr := func() (int, error) {
    return http.StatusInternalServerError, errors.New("unable to leave study group")
  }

  err := server.DB.Get(
    &studyGroup,
    "SELECT user_id, members, waitlist, available_spots FROM study_groups WHERE id = $1",
    studyGroupID,
  )

  switch {
    case err == sql.ErrNoRows:
      return http.StatusNotFound, errors.New("study group not found")
    case err != nil:
      return internalErr()
  }

  err = server.DB.Get(
    &user,
    "SELECT study_groups, waitlists FROM users WHERE id = $1",
    userID,
  )

  switch {
    case err == sql.ErrNoRows:
      return http.StatusNotFound, errors.New("user not found")
    case err != nil:
      return internalErr()
  }

  if err = studyGroup.AddUserToWaitlist(userID); err != nil {
    return http.StatusForbidden, err
  }
  if err = user.AddStudyGroupToWaitlists(studyGroupID); err != nil {
    return http.StatusForbidden, err
  }

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

    if studyGroup.Waitlist.String == "" {
      _, err = tx.Exec(
        "UPDATE study_groups SET waitlist = null, available_spots = $1 WHERE id = $2",
        studyGroup.AvailableSpots,
        studyGroupID,
      )
    } else {
      _, err = tx.Exec(
        "UPDATE study_groups SET waitlist = $1, available_spots = $2 WHERE id = $3",
        studyGroup.Waitlist.String,
        studyGroup.AvailableSpots,
        studyGroupID,
      )
    }

    if user.Waitlists.String == "" {
      _, err = tx.Exec("UPDATE users SET waitlists = null WHERE id = $1",
        userID,
      )
    } else {
      _, err = tx.Exec("UPDATE users SET waitlists = $1 WHERE id = $2",
        user.Waitlists.String,
        userID,
      )
    }


    err = tx.Commit()
  }

  return http.StatusOK, nil
}

func LeaveStudyGroup(studyGroupID, userID string) (int, error) {
  var user models.User
  var studyGroup models.StudyGroup

  internalErr := func() (int, error) {
    return http.StatusInternalServerError, errors.New("unable to leave study group")
  }

  err := server.DB.Get(
    &studyGroup,
    "SELECT user_id, members, waitlist, available_spots FROM study_groups WHERE id = $1",
    studyGroupID,
  )

  switch {
    case err == sql.ErrNoRows:
      return http.StatusNotFound, errors.New("study group not found")
    case err != nil:
      return internalErr()
  }

  err = server.DB.Get(
    &user,
    "SELECT study_groups, waitlists FROM users WHERE id = $1",
    userID,
  )

  switch {
    case err == sql.ErrNoRows:
      return http.StatusNotFound, errors.New("user not found")
    case err != nil:
      return internalErr()
  }

  sgColumnName, sgColumnVal, err := studyGroup.RemoveUser(userID)
  if err != nil { return http.StatusForbidden, err }

  uColumnName, uColumnVal, err := user.LeaveStudyGroup(studyGroupID)
  if err != nil { return http.StatusForbidden, err }

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

    // TODO: refactor dynamic query
    if sgColumnVal.String == "" {
      _, err = tx.Exec(
       "UPDATE study_groups SET " + sgColumnName + " = null, available_spots = $1 WHERE id = $2",
        studyGroup.AvailableSpots,
        studyGroupID,
      )
    } else {
      _, err = tx.Exec(
       "UPDATE study_groups SET " + sgColumnName + " = $1, available_spots = $2 WHERE id = $3",
        sgColumnVal.String,
        studyGroup.AvailableSpots,
        studyGroupID,
      )
    }

    if uColumnVal.String == "" {
      _, err = tx.Exec("UPDATE users SET " + uColumnName + " = null WHERE id = $1",
        userID,
      )
    } else {
      _, err = tx.Exec("UPDATE users SET " + uColumnName + " = $1 WHERE id = $2",
        uColumnVal.String,
        userID,
      )
    }


    err = tx.Commit()
  }

  return http.StatusOK, nil
}
