package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/prosperoa/study-groups/src/models"
	"github.com/prosperoa/study-groups/src/server"
	"golang.org/x/crypto/bcrypt"
)

func Login(credentials models.LoginCredentials) (models.User, int, error) {
	var user models.User

	err := server.DB.Get(&user, "SELECT * FROM users WHERE email = $1", credentials.Email)

	switch {
	case err == sql.ErrNoRows:
		return user, http.StatusUnauthorized, errors.New("account does not exist")
	case err != nil:
		return user, http.StatusInternalServerError, errors.New("unable to login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return user, http.StatusBadRequest, errors.New("incorrect password")
	}

	return user, http.StatusOK, nil
}

func Signup(credentials models.SignUpCredentials) (models.User, int, error) {
	var user models.User
	var accountExists bool
	errMsg := "unable to create account"

	err := server.DB.Get(&accountExists, "SELECT exists(SELECT 1 FROM users WHERE email = $1)",
		credentials.Email,
	)

	if accountExists {
		return user, http.StatusForbidden, errors.New("account already exists")
	} else if err != nil {
		return user, http.StatusInternalServerError, errors.New(errMsg)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.MinCost)

	if err != nil {
		return user, http.StatusInternalServerError, errors.New(errMsg)
	}

	err = server.DB.Get(
		&user,
		`INSERT INTO users (first_name, last_name, email, password, created_on, updated_on)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING *`,
		credentials.FirstName, credentials.LastName, credentials.Email, passwordHash, time.Now(), time.Now(),
	)

	if err != nil {
		return user, http.StatusInternalServerError, errors.New(errMsg)
	}

	// add user to mailchimp list
	if err = mailchimp.SetKey(os.Getenv("MAILCHIMP_API_KEY")); err != nil {
		log.Println(err.Error())
	}

	params := &members.NewParams{
		EmailAddress: user.Email,
		Status:       members.StatusSubscribed,
	}

	if _, err = members.New("4d6392ba4d", params); err != nil {
		log.Println(err.Error())
	}

	return user, http.StatusOK, nil
}
