package models

import (
  "errors"
  "strings"

  "github.com/prosperoa/study-groups/src/utils"
  "gopkg.in/guregu/null.v3"
)

type User struct {
  ID          int         `db:"id"           json:"id"`
  FirstName   string      `db:"first_name"   json:"first_name"`
  LastName    null.String `db:"last_name"    json:"last_name"`
  Email       string      `db:"email"        json:"email"`
  Avatar      null.String `db:"avatar"       json:"avatar"`
  Bio         null.String `db:"bio"          json:"bio"`
  School      null.String `db:"school"       json:"school"`
  Major1      null.String `db:"major1"       json:"major1"`
  Major2      null.String `db:"major2"       json:"major2"`
  Minor       null.String `db:"minor"        json:"minor"`
  Courses     null.String `db:"courses"      json:"courses"`
  StudyGroups null.String `db:"study_groups" json:"-"`
  Password    string      `db:"password"     json:"-"`
}

func (u *User) LeaveStudyGroup(studyGroupID string) error {
  studyGroups := strings.Split(u.StudyGroups.String, ",")

  if !utils.Contains(studyGroups, studyGroupID) {
    return errors.New("user is not a member of study group")
  }

  if len(studyGroups) > 1 {
    u.StudyGroups = null.StringFrom(strings.Join(utils.Splice(studyGroups, studyGroupID), ","))
  } else {
    u.StudyGroups = null.String{}
  }

  return nil
}
