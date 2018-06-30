package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prosperoa/study-groups/src/controllers"
	"github.com/prosperoa/study-groups/src/models"
	"github.com/prosperoa/study-groups/src/server"
	"github.com/prosperoa/study-groups/src/utils"
)

func GetStudyGroup(c *gin.Context) {
	id := c.Param("id")

	if !utils.IsInt(id) {
		server.Respond(c, nil, "invalid study group id", http.StatusBadRequest)
		return
	}

	studyGroup, status, err := controllers.GetStudyGroup(id)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroup, "", status)
}

func GetStudyGroups(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pageSize := c.DefaultQuery("page_size", "30")

	if !utils.IsInt(page) || !utils.IsInt(pageSize) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	studyGroups, status, err := controllers.GetStudyGroups(p, ps)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroups, "", status)
}

func JoinStudyGroup(c *gin.Context) {
	var userID models.UserID
	studyGroupID := c.Param("id")

	if err := c.ShouldBindWith(&userID, binding.JSON); err != nil {
		server.Respond(c, nil, "missing user id", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(userID); err != nil || !utils.IsInt(studyGroupID) {
		server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
		return
	}

	status, err := controllers.JoinStudyGroup(studyGroupID, userID.String())

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, nil, "user added to study group waitlist", status)
}

func DeleteStudyGroup(c *gin.Context) {
	var userID models.UserID
	studyGroupID := c.Param("id")

	if err := c.ShouldBindWith(&userID, binding.JSON); err != nil {
		server.Respond(c, nil, "missing user id", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(userID); err != nil || !utils.IsInt(studyGroupID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	status, err := controllers.DeleteStudyGroup(studyGroupID, userID.String())

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, nil, "study group successfully deleted", status)
}

func LeaveStudyGroup(c *gin.Context) {
	var userID models.UserID
	studyGroupID := c.Param("id")

	if err := c.ShouldBindWith(&userID, binding.JSON); err != nil {
		server.Respond(c, nil, "missing user id", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(userID); err != nil || !utils.IsInt(studyGroupID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	status, err := controllers.LeaveStudyGroup(studyGroupID, userID.String())

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, nil, "user removed from study group", status)
}
