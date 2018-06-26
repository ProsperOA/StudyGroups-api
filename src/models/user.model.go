package models

import (
  "errors"
  "reflect"
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
  Waitlists   null.String `db:"waitlists"    json:"-"`
  Password    string      `db:"password"     json:"-"`
  CreatedOn   string      `db:"created_on"   json:"-"`
  UpdatedOn   string      `db:"updated_on"   json:"-"`
}

func (u *User) AddStudyGroupToWaitlists(studyGroupID string) error {
  studyGroups := strings.Split(u.StudyGroups.String, ",")
  waitlists := strings.Split(u.Waitlists.String, ",")

  if utils.Contains(waitlists, studyGroupID) {
    return errors.New("user is already waitlisted")
  } else if utils.Contains(studyGroups, studyGroupID) {
    return errors.New("user is already in study group")
  }

  if u.Waitlists.String == "" {
    u.Waitlists = null.StringFrom(studyGroupID)
  } else {
    waitlists = append(waitlists, studyGroupID)
    u.Waitlists = null.StringFrom(strings.Join(waitlists, ","))
  }

  return nil
}

func (u *User) LeaveStudyGroup(studyGroupID string) (string, null.String, error) {
  studyGroups := strings.Split(u.StudyGroups.String, ",")
  waitlists := strings.Split(u.Waitlists.String, ",")

  if !utils.Contains(studyGroups, studyGroupID) && !utils.Contains(waitlists, studyGroupID) {
    return "", null.String{}, errors.New("user is not waitlisted or a member of study group")
  }

  var uField *null.String
  var field *[]string
  var columnName string

  if utils.Contains(studyGroups, studyGroupID) {
    uField = &u.StudyGroups
    field = &studyGroups

    f, _ := reflect.TypeOf(u).Elem().FieldByName("StudyGroups")
    columnName = f.Tag.Get("db")
  } else {
    uField = &u.Waitlists
    field = &waitlists

    f, _ := reflect.TypeOf(u).Elem().FieldByName("Waitlists")
    columnName = f.Tag.Get("db")
  }

  if len(*field) > 1 {
    *uField = null.StringFrom(strings.Join(utils.Splice(*field, studyGroupID), ","))
  } else {
    *uField = null.String{}
  }

  return columnName, *uField, nil
}
