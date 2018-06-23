package models

import "gopkg.in/guregu/null.v3"

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
