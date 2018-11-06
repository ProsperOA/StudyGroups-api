package models

import (
	// "database/sql"
	"errors"
	// "fmt"
	// "log"
	// "mime/multipart"
	// "net/http"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/prosperoa/study-groups/src/server"
	"github.com/prosperoa/study-groups/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v3"
)

type User struct {
	ID          int                `db:"id"           json:"id"`
	FirstName   string             `db:"first_name"   json:"first_name"`
	LastName    null.String        `db:"last_name"    json:"last_name"`
	Email       string             `db:"email"        json:"email"`
	Avatar      null.String        `db:"avatar"       json:"avatar"`
	Bio         null.String        `db:"bio"          json:"bio"`
	School      null.String        `db:"school"       json:"school"`
	Major1      null.String        `db:"major1"       json:"major1"`
	Major2      null.String        `db:"major2"       json:"major2"`
	Minor       null.String        `db:"minor"        json:"minor"`
	Courses     types.NullJSONText `db:"courses"      json:"courses"`
	StudyGroups null.String        `db:"study_groups" json:"-"`
	Waitlists   null.String        `db:"waitlists"    json:"-"`
	Password    string             `db:"password"     json:"-"`
	CreatedOn   string             `db:"created_on"   json:"-"`
	UpdatedOn   string             `db:"updated_on"   json:"-"`
}

type Users []User

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

func (u *User) Get() error {
	if u.ID == 0 { return errors.New("missing user ID") }

	return server.DB.Get(u, "SELECT * FROM users WHERE id = $1", u.ID)
}

func (u *Users) Get(page, pageSize int) error {
  if page < 0 { page = 0 }
  if pageSize < 0 { pageSize = 30 }

	return server.DB.Select(u, "SELECT * FROM users LIMIT $1 OFFSET $2",
		pageSize, pageSize * page,
	)
}

func (u *User) Delete() error {
  var passwordHash string

  if u.ID == 0 || u.Password == "" {
    return errors.New("invalid password")
  }

  err := server.DB.Get(&passwordHash, "SELECT password FROM users WHERE id = $1", u.ID)
	if err != nil {	return err }

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(u.Password))
	if err != nil {	return err }

  err = server.DB.Get(u, "DELETE FROM users WHERE id = $1 RETURNING email, avatar", u.ID)
	if err != nil {	return err }

	return nil
}

// func GetUserStudyGroups(userID string, page, pageSize int) ([]models.StudyGroup, int, error) {
// 	var studyGroups []models.StudyGroup

// 	err := server.DB.Select(
// 		&studyGroups,
// 		"SELECT * FROM study_groups WHERE user_id = $1 ORDER BY updated_on DESC LIMIT $2 OFFSET $3",
// 		userID,
// 		pageSize,
// 		pageSize*page,
// 	)

// 	switch {
// 	case err == sql.ErrNoRows, len(studyGroups) == 0:
// 		return studyGroups, http.StatusNotFound, errors.New("no users study groups found")
// 	case err != nil:
// 		return studyGroups, http.StatusInternalServerError, errors.New("unable to get user's study groups")
// 	}

// 	return studyGroups, http.StatusOK, nil
// }

func (u *User) SetAvatar(avatarURL string) error {
  if u.ID == 0 || u.Avatar.String == "" {
    return errors.New("missing user id or avatar url")
  }

  _, err := server.DB.Exec(
   "UPDATE users SET avatar = $1 WHERE id = $2",
    avatarURL,
    u.ID,
  )

	return err
}

func (u *User) UpdateAccount() error {
  if u.ID == 0 || u.FirstName == "" {
    return errors.New("user id and first name required")
  }

	return server.DB.Get(
    u,
   `UPDATE
      users
    SET
      first_name = $1,
      last_name  = $2,
      bio        = $3,
      school     = $4,
      major1     = $5,
      major2     = $6,
      minor      = $7,
      updated_on = $8
    WHERE
      id = $9
    RETURNING
      *`,
		u.FirstName,
		u.LastName,
		u.Bio,
		u.School,
		u.Major1,
		u.Major2,
		u.Minor,
		time.Now(),
		u.ID,
	)
}

func (u *User) ChangePassword(newPassword string) error {
	var currentPasswordHash string

	err := server.DB.Get(
    &currentPasswordHash,
    "SELECT password FROM users WHERE id = $1",
		u.ID,
  )
  if err != nil { return errors.New("user not found") }

	err = bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(u.Password))
	if err != nil {	return err }

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil { return err }

	_, err = server.DB.Exec(
   "UPDATE users SET password = $1 WHERE id = $2",
    newPasswordHash,
    u.ID,
  )

  return err
}

func (u *User) UpdateCourses () error {
  if u.ID == 0 || u.Courses.JSONText == nil {
    return errors.New("invalid user id or courses")
  }

	_, err := server.DB.Exec(
   "UPDATE users SET courses = $1 WHERE id = $2",
    u.Courses,
    u.ID,
  )

  return err
}
