package controllers

import (
  "database/sql"
  "errors"
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
