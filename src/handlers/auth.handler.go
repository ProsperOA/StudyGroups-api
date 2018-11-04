package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prosperoa/study-groups/src/controllers"
	"github.com/prosperoa/study-groups/src/email-notifications"
	"github.com/prosperoa/study-groups/src/models"
	"github.com/prosperoa/study-groups/src/server"
)

func Login(c *gin.Context) {
	var credentials models.LoginCredentials

	if err := c.ShouldBindWith(&credentials, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(credentials); err != nil {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	user, status, err := controllers.Login(credentials)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	authToken, err := server.GenerateAuthToken(strconv.Itoa(user.ID))
	if err != nil {
		server.Respond(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"auth_token": authToken,
		"user":       user,
	}

	server.Respond(c, data, "", status)
}

func Signup(c *gin.Context) {
	var credentials models.SignUpCredentials

	if err := c.ShouldBindWith(&credentials, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(credentials); err != nil {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	err := server.Validate.VarWithValue(credentials.Password,
		credentials.ConfirmPassword, "eqfield",
	)
	if err != nil {
		server.Respond(c, nil, "passwords must match", http.StatusBadRequest)
		return
	}

	if err := server.ValidateEmail(credentials.Email); err != nil {
		server.Respond(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	user, status, err := controllers.Signup(credentials)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	authToken, err := server.GenerateAuthToken(strconv.Itoa(user.ID))
	if err != nil {
		server.Respond(c, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = emails.NewUserNotification(user.FirstName, user.Email); err != nil {
		log.Println(err.Error())
	}

	data := map[string]interface{}{
		"auth_token": authToken,
		"user":       user,
	}

	server.Respond(c, data, "", status)
}
