package models

import "strconv"

type UserID struct {
	Value int `json:"user_id" validate:"required,gt=0"`
}

type LoginCredentials struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type SignUpCredentials struct {
	LoginCredentials
	FirstName       string `json:"first_name"       validate:"required,min=1,max=20"`
	LastName        string `json:"last_name"        validate:"max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=50"`
}

type Account struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=20"`
	LastName  string `json:"last_name"  validate:"max=20"`
	Bio       string `json:"bio"        validate:"max=280"`
	School    string `json:"school"     validate:"max=20"`
	Major1    string `json:"major1"     validate:"max=40"`
	Major2    string `json:"major2"     validate:"max=40"`
	Minor     string `json:"minor"      validate:"max=40"`
}

type Password struct {
	Value string `json:"password" validate:"required,min=6,max=50"`
}

type ChangePassword struct {
	New     string `json:"new_password"     validate:"required,min=6,max=50,excludesall= "`
	Confirm string `json:"confirm_password" validate:"required,min=6,max=50,excludesall= "`
	Current string `json:"current_password" validate:"required"`
}

type StudyGroupsFilter struct {
	BaseFilter
	AvailableSpots int    `json:"available_spots" validate:"min=1"`
	Location       string `json:"location"`
	MeetingDate    string `json:"meeting_date"`
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	Instructor     string `json:"instructor"`
	Term           string `json:"term"`
}

func (u UserID) String() string {
	return strconv.Itoa(u.Value)
}
