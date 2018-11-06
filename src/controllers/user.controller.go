package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	 mailchimp "github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prosperoa/study-groups/src/models"
	"github.com/prosperoa/study-groups/src/server"
	"github.com/prosperoa/study-groups/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v3"
)

func GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil || userID == 0 {
		server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
		return
	}

	user := models.User{ID: userID}
	err = user.Get()

	switch {
		case err == sql.ErrNoRows:
			server.Respond(c, nil, "user not found", http.StatusNotFound)
			return
		case err != nil:
			server.Respond(c, nil, "unable to get user", http.StatusInternalServerError)
			return
	}

	server.Respond(c, user, "", http.StatusOK)
}

func GetUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pageSize := c.DefaultQuery("page_size", "30")

	if page == "" || pageSize == "" {
		server.Respond(c, nil, "missing page or page size", http.StatusBadRequest)
		return
	}

	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	var users models.Users
	err := users.Get(p, ps)

	switch {
		case err == sql.ErrNoRows:
			server.Respond(c, nil, "no users found", http.StatusNotFound)
			return
		case err != nil:
			server.Respond(c, nil, "unable to get users", http.StatusInternalServerError)
			return
	}

	server.Respond(c, users, "", http.StatusOK)
}

func DeleteUser(c *gin.Context) {
	var password models.Password
	userID, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindWith(&password, binding.JSON); err != nil || userID == 0 {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       userID,
		Password: password.Value,
	}

	err := user.Delete()

	switch {
		case err == sql.ErrNoRows:
			server.Respond(c, nil, "account doesn't exist", http.StatusNotFound)
			return
		case err == bcrypt.ErrMismatchedHashAndPassword:
			server.Respond(c, nil, "incorrect password", http.StatusForbidden)
			return
		case err != nil:
			server.Respond(c, nil, "unable to delete account", http.StatusInternalServerError)
			return
	}

	// delete avatar if it's not the stock avatar image
	if !strings.Contains(user.Avatar.String, "stock-avatar") {
		_, err = server.S3Service.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(server.S3Bucket),
			Key:    aws.String(strings.TrimPrefix(user.Avatar.String, server.S3BucketURL)),
		})

		if err != nil { log.Println(err.Error()) }
	}

	// remove user from mailchimp list
	var mailchimpID string

	params := &members.GetParams{Status: members.StatusSubscribed}
	listMembers, err := members.Get("4d6392ba4d", params)

	if err == nil {
		for _, v := range listMembers.Members {
			if v.EmailAddress == user.Email {
				mailchimpID = v.ID
				break
			}
		}

		if err = mailchimp.SetKey(os.Getenv("MAILCHIMP_API_KEY")); err == nil {
			if err = members.Delete("4d6392ba4d", mailchimpID); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}

	server.Respond(c, nil, "account successfully deleted", http.StatusOK)
}

// func GetUserStudyGroups(c *gin.Context) {
// 	userID := c.Param("id")
// 	page := c.DefaultQuery("page", "0")
// 	pageSize := c.DefaultQuery("page_size", "30")

// 	if !utils.IsInt(userID) || !utils.IsInt(page) || !utils.IsInt(pageSize) {
// 		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
// 		return
// 	}

// 	p, _ := strconv.Atoi(page)
// 	ps, _ := strconv.Atoi(pageSize)

// 	studyGroups, status, err := controllers.GetUserStudyGroups(userID, p, ps)

// 	if err != nil {
// 		server.Respond(c, nil, err.Error(), status)
// 		return
// 	}

// 	server.Respond(c, studyGroups, "", status)
// }

func UploadAvatar(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	file, err := c.FormFile("image")
	errMsg := "unable to upload avatar"

	switch {
		case userID == 0:
			server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
			return
		case err != nil:
			server.Respond(c, nil, "invalid image", http.StatusBadRequest)
			return
		case file.Size / utils.MB > utils.MB * 2:
			server.Respond(c, nil, "image size must be 2MB or less", http.StatusBadRequest)
			return
	}

	user := models.User{ID: userID}

	if err = user.Get(); err != nil {
		server.Respond(c, nil, errMsg, http.StatusInternalServerError)
		return
	}

	// delete old avatar if it's not the stock avatar
	if !strings.Contains(user.Avatar.String, "stock-avatar") {
		_, err = server.S3Service.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(server.S3Bucket),
			Key:    aws.String(strings.TrimPrefix(user.Avatar.String, server.S3BucketURL)),
		})

		if err != nil {	log.Println(err.Error()) }
	}

	// construct image filename and upload
	ext := filepath.Ext(file.Filename)
	image, _ := file.Open()
	newAvatarFilename := fmt.Sprintf("%d-%s", user.ID, utils.RandString(16)+ext)

	result, err := server.S3Uploader.Upload(&s3manager.UploadInput{
		Body:   image,
		Bucket: aws.String(server.S3Bucket),
		Key:    aws.String("images/user-avatars/" + newAvatarFilename),
		ACL:    aws.String("public-read"),
	})

	newAvatarURL := result.Location

	if err != nil || user.SetAvatar(newAvatarURL) != nil {
		server.Respond(c, nil, errMsg, http.StatusInternalServerError)
		return
	}

	server.Respond(c, newAvatarURL, "", http.StatusOK)
}

func UpdateAccount(c *gin.Context) {
	var account models.Account
	userID, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindWith(&account, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	accountValues := reflect.ValueOf(&account).Elem()

	// remove extraneous whitespace from strings
	for i := 0; i < accountValues.NumField(); i++ {
		field := accountValues.Field(i)
		if field.Interface() != "string" {
			continue
		}

		str := field.Interface().(string)
		field.SetString(utils.Trim(str))
	}

	if err := server.Validate.Struct(account); err != nil || userID == 0 {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:        userID,
		FirstName: account.FirstName,
		LastName:  null.NewString(account.LastName, account.LastName != ""),
		Bio:       null.NewString(account.Bio, account.Bio != ""),
		School:    null.NewString(account.School, account.School != ""),
		Major1:    null.NewString(account.Major1, account.Major1 != ""),
		Major2:    null.NewString(account.Major2, account.Major2 != ""),
		Minor:     null.NewString(account.Minor, account.Minor != ""),
	}

	if err := user.UpdateAccount(); err != nil {
		server.Respond(c, nil, "unable to update account", http.StatusBadRequest)
		return
	}

	server.Respond(c, user, "", http.StatusOK)
}

func ChangePassword(c *gin.Context) {
	var passwords models.ChangePassword
	userID, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindWith(&passwords, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(passwords); err != nil || userID == 0 {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	err := server.Validate.VarWithValue(passwords.New, passwords.Confirm, "eqfield")
	if err != nil {
		server.Respond(c, nil, "passwords must match", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       userID,
		Password: passwords.Current,
	}

	err = user.ChangePassword(passwords.New)

	switch {
		case err == sql.ErrNoRows:
			server.Respond(c, nil, err.Error(), http.StatusNotFound)
			return
		case err == bcrypt.ErrMismatchedHashAndPassword:
			server.Respond(c, nil, "incorrect password", http.StatusUnauthorized)
			return
		case err != nil:
			server.Respond(c, nil, err.Error(), http.StatusInternalServerError)
			return
	}

	server.Respond(c, nil, "password successfully changed", http.StatusOK)
}

func UpdateCourses(c *gin.Context) {
	var coursesJSON []models.Course
	userID, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindWith(&coursesJSON, binding.JSON); err != nil || userID == 0 {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	courses, err := json.Marshal(coursesJSON)
	if err != nil {
		server.Respond(c, nil, "invalid courses JSON format", http.StatusBadRequest)
		return
	}

	user := models.User{ID: userID}
	user.Courses.Scan(courses)

	if err = user.UpdateCourses(); err != nil {
		server.Respond(c, nil, "unable to update courses", http.StatusInternalServerError)
		return
	}

	server.Respond(c, nil, "courses successfully updated", http.StatusOK)
}
