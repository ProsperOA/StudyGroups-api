package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prosperoa/study-groups/src/controllers"
	"github.com/prosperoa/study-groups/src/models"
	"github.com/prosperoa/study-groups/src/server"
	"github.com/prosperoa/study-groups/src/utils"
)

func GetUser(c *gin.Context) {
	userID := c.Param("id")

	if !utils.IsInt(userID) {
		server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
		return
	}

	user, status, err := controllers.GetUser(userID)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, user, "", status)
}

func GetUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pageSize := c.DefaultQuery("page_size", "30")

	if !utils.IsInt(page) || !utils.IsInt(pageSize) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	users, status, err := controllers.GetUsers(p, ps)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, users, "", status)
}

func DeleteUser(c *gin.Context) {
	var password models.Password
	userID := c.Param("id")

	if err := c.ShouldBindWith(&password, binding.JSON); err != nil || !utils.IsInt(userID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	status, err := controllers.DeleteUser(userID, password.Value)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, nil, "account successfully deleted", status)
}

func GetUserStudyGroups(c *gin.Context) {
	userID := c.Param("id")
	page := c.DefaultQuery("page", "0")
	pageSize := c.DefaultQuery("page_size", "30")

	if !utils.IsInt(userID) || !utils.IsInt(page) || !utils.IsInt(pageSize) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	studyGroups, status, err := controllers.GetUserStudyGroups(userID, p, ps)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroups, "", status)
}

func UploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	file, err := c.FormFile("image")

	if !utils.IsInt(userID) {
		server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
		return
	}

	if err != nil {
		server.Respond(c, nil, "unable to get image", http.StatusBadRequest)
		return
	}

	if file.Size/utils.MB > utils.MB*2 {
		server.Respond(c, nil, "image size must be 2MB or less", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(file.Filename)
	image, _ := file.Open()
	avatarURL, status, err := controllers.UploadAvatar(userID, ext, image)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, avatarURL, "", status)
}

func UpdateAccount(c *gin.Context) {
	var account models.Account
	userID := c.Param("id")

	if err := c.ShouldBindWith(&account, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	valuesPtr := reflect.ValueOf(&account)
	values := valuesPtr.Elem()

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		if field.Interface() != "string" {
			continue
		}

		str := field.Interface().(string)
		field.SetString(utils.Trim(str))
	}

	if err := server.Validate.Struct(account); err != nil || !utils.IsInt(userID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	user, status, err := controllers.UpdateAccount(userID, account)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, user, "", status)
}

func ChangePassword(c *gin.Context) {
	var passwords models.ChangePassword
	userID := c.Param("id")

	if err := c.ShouldBindWith(&passwords, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(passwords); err != nil || !utils.IsInt(userID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	err := server.Validate.VarWithValue(passwords.New, passwords.Confirm, "eqfield")
	if err != nil {
		server.Respond(c, nil, "passwords must match", http.StatusBadRequest)
		return
	}

	user, status, err := controllers.ChangePassword(userID, passwords)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, user, "", status)
}

func UpdateCourses (c *gin.Context) {
	var courses []models.Course
	userID := c.Param("id")

	if err := c.ShouldBindWith(&courses, binding.JSON); err != nil || !utils.IsInt(userID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	coursesB, err := json.Marshal(courses)
	if err != nil {
		server.Respond(c, nil, "invalid courses format", http.StatusBadRequest)
		return
	}

	status, err := controllers.UpdateCourses(userID, coursesB)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, nil, "courses successfully updated", status)
}
