package models

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx/types"
	"github.com/prosperoa/study-groups/src/utils"
	"gopkg.in/guregu/null.v3"
)

type StudyGroup struct {
	ID             int                 `db:"id"              json:"id"`
	UserID         int                 `db:"user_id"         json:"user_id"`
	Name           string              `db:"name"            json:"name"`
	Members        null.String         `db:"members"         json:"members"`
	MembersLimit   null.Int            `db:"members_limit"   json:"members_limit"`
	AvailableSpots int                 `db:"available_spots" json:"available_spots"`
	Location       null.String         `db:"location"        json:"location"`
	Description    null.String         `db:"description"     json:"description"`
	MeetingDate    null.String         `db:"meeting_date"    json:"meeting_date"`
	Course         types.NullJSONText  `db:"course"          json:"course"`
	Waitlist       null.String         `db:"waitlist"        json:"waitlist"`
	CreatedAt      string              `db:"created_on"      json:"-"`
	UpdatedAt      string              `db:"updated_on"      json:"-"`
}

func (sg *StudyGroup) AddUserToWaitlist(userID string) error {
	if sg.AvailableSpots == 0 {
		return errors.New("study group members limit reached")
	}

	if utils.Contains(strings.Split(sg.Members.String, ","), userID) {
		return errors.New("user is already in study group")
	}

	uID, _ := strconv.Atoi(userID)
	if sg.UserID == uID {
		return errors.New("user is owner of group")
	}

	waitlist := strings.Split(sg.Waitlist.String, ",")
	if utils.Contains(waitlist, userID) {
		return errors.New("user is already a member of study group")
	}

	if sg.Waitlist.String == "" {
		sg.Waitlist = null.StringFrom(userID)
	} else {
		waitlist = append(waitlist, userID)
		sg.Waitlist = null.StringFrom(strings.Join(waitlist, ","))
	}

	sg.AvailableSpots -= 1

	return nil
}

func (sg *StudyGroup) RemoveUser(userID string) (string, null.String, error) {
	members := strings.Split(sg.Members.String, ",")
	waitlist := strings.Split(sg.Waitlist.String, ",")

	uID, _ := strconv.Atoi(userID)
	if sg.UserID == uID {
		return "", null.String{}, errors.New("user is owner of study group")
	}

	if !utils.Contains(members, userID) && !utils.Contains(waitlist, userID) {
		return "", null.String{}, errors.New("user is not waitlisted or a member of study group")
	}

	var sgField *null.String
	var field *[]string
	var columnName string

	if utils.Contains(members, userID) {
		sgField = &sg.Members
		field = &members

		f, _ := reflect.TypeOf(sg).Elem().FieldByName("Members")
		columnName = f.Tag.Get("db")
	} else {
		sgField = &sg.Waitlist
		field = &waitlist

		f, _ := reflect.TypeOf(sg).Elem().FieldByName("Waitlist")
		columnName = f.Tag.Get("db")
	}

	if len(*field) > 1 {
		*sgField = null.StringFrom(strings.Join(utils.Splice(*field, userID), ","))
	} else {
		*sgField = null.String{}
	}

	sg.AvailableSpots += 1

	return columnName, *sgField, nil
}
