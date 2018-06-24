package models

import (
  "errors"
  "strings"

  "github.com/prosperoa/study-groups/src/utils"
  "gopkg.in/guregu/null.v3"
)

type StudyGroup struct {
  ID             int         `db:"id"              json:"id"`
  UserID         int         `db:"user_id"         json:"user_id"`
  Name           string      `db:"name"            json:"name"`
  Members        null.String `db:"members"         json:"members"`
  MembersLimit   null.Int    `db:"members_limit"   json:"members_limit"`
  AvailableSpots null.Int    `db:"available_spots" json:"available_spots"`
  Location       null.String `db:"location"        json:"location"`
  Description    null.String `db:"description"     json:"description"`
  MeetingDate    null.String `db:"meeting_date"    json:"meeting_date"`
  Course         null.String `db:"course"          json:"course"`
  CreatedAt      string      `db:"created_at"      json:"created_at"`
  UpdatedAt      string      `db:"updated_at"      json:"updated_at"`
}

func (sg *StudyGroup) RemoveMember(userID string) error {
  members := strings.Split(sg.Members.String, ",")

  if !utils.Contains(members, userID) {
    return errors.New("user is not a member of " + sg.Name)
  }

  if len(members) > 1 {
    sg.Members = null.StringFrom(strings.Join(utils.Splice(members, userID), ","))
  } else {
    sg.Members = null.String{}
  }

  return nil
}
