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
	uID, _            := strconv.Atoi(c.Query("user_id"))
	userID            := models.UserID{Value: uID}

	pageIndex, _      := strconv.Atoi(c.DefaultQuery("page_index", "0"))
	pageSize, _       := strconv.Atoi(c.DefaultQuery("page_size", "30"))
	availableSpots, _ := strconv.Atoi(c.DefaultQuery("available_spots", "1"))

	filter := models.StudyGroupsFilter{
		BaseFilter: models.BaseFilter{
			PageIndex: pageIndex,
			PageSize:  pageSize,
		},
		AvailableSpots: availableSpots,
		StudyGroupName: c.Query("study_group_name"),
		Location:       c.Query("location"),
		MeetingDate:    c.Query("meeting_date"),
		CourseCode:     c.Query("course_code"),
		CourseName:     c.Query("course_name"),
		Instructor:     c.Query("instructor"),
		Term:           c.Query("term"),
	}

	if err := server.Validate.Struct(userID); err != nil {
		server.Respond(c, nil, "invalid user id", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(filter); err != nil {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	studyGroups, status, err := controllers.GetStudyGroups(filter, userID.Value)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	var message string
	if len(studyGroups) == 0 { message = "no study groups found" }

	server.Respond(c, studyGroups, message, status)
}

func GetStudyGroupMembers(c *gin.Context) {
	studyGroupID := c.Param("id")

	if !utils.IsInt(studyGroupID) {
		server.Respond(c, nil, "invalid study group id", http.StatusBadRequest)
		return
	}

	studyGroupMembers, status, err := controllers.GetStudyGroupMembers(studyGroupID)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroupMembers, "", status)
}

func CreateStudyGroup(c *gin.Context) {
	var studyGroup models.StudyGroup

	if err := c.ShouldBindWith(&studyGroup, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(studyGroup); err != nil {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	studyGroup, status, err := controllers.CreateStudyGroup(studyGroup)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroup, "study group created", status)

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

	studyGroup, status, err := controllers.JoinStudyGroup(studyGroupID, userID.String())

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroup, "user added to study group waitlist", status)
}

func MoveUserFromWaitlistToMembers(c *gin.Context) {
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

	studyGroup, status, err := controllers.MoveUserFromWaitlistToMembers(studyGroupID, userID.String())

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, studyGroup, "", status)
}

func UpdateStudyGroup(c *gin.Context) {
	var studyGroup models.StudyGroup
	studyGroupID := c.Param("id")

	if err := c.ShouldBindWith(&studyGroup, binding.JSON); err != nil {
		server.Respond(c, nil, "missing params", http.StatusBadRequest)
		return
	}

	if err := server.Validate.Struct(studyGroup); err != nil || !utils.IsInt(studyGroupID) {
		server.Respond(c, nil, "invalid params", http.StatusBadRequest)
		return
	}

	updatedStudyGroup, status, err := controllers.UpdateStudyGroup(studyGroup)

	if err != nil {
		server.Respond(c, nil, err.Error(), status)
		return
	}

	server.Respond(c, updatedStudyGroup, "", status)
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
