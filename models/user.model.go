package models

import "gopkg.in/guregu/null.v3"

type User struct {
  ID        int         `db:"id"         json:"id"`
  FirstName string      `db:"first_name" json:"first_name"`
  LastName  null.String `db:"last_name"  json:"last_name"`
  Email     string      `db:"email"      json:"email"`
  Avatar    null.String `db:"avatar"     json:"avatar"`
  Bio       null.String `db:"bio"        json:"bio"`
  School    null.String `db:"school"     json:"school"`
  Major1    null.String `db:"major1"     json:"major1"`
  Major2    null.String `db:"major2"     json:"major2"`
  Minor     null.String `db:"minor"      json:"minor"`
  Courses   null.String `db:"courses"    json:"courses"`
  Password  string      `db:"password"   json:"-"`
}
