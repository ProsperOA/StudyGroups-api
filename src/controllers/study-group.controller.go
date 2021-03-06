package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prosperoa/study-groups/src/models"
	"github.com/prosperoa/study-groups/src/server"
)

func GetStudyGroup(id string) (models.StudyGroup, int, error) {
	var studyGroup models.StudyGroup

	err := server.DB.Get(&studyGroup, "SELECT * FROM study_groups WHERE id = $1", id)

	switch {
	case err == sql.ErrNoRows:
		return studyGroup, http.StatusNotFound, errors.New("study group not found")
	case err != nil:
		return studyGroup, http.StatusInternalServerError, errors.New(
			"unable to get study group",
		)
	}

	return studyGroup, http.StatusOK, nil
}

func GetStudyGroups(filter models.StudyGroupsFilter, userID int) ([]models.StudyGroup, int, error) {
	var studyGroups []models.StudyGroup

	query := fmt.Sprintf(
		"SELECT * FROM study_groups WHERE user_id != %d AND available_spots >= %d",
		userID,
		filter.AvailableSpots,
	)

	query += fmt.Sprintf(`
		AND %d != ANY(
			(SELECT	regexp_split_to_array(members, ',')	FROM study_groups)::int[]
		)`,
		userID,
	)

	if filter.StudyGroupName != "" {
		query += fmt.Sprintf(" AND levenshtein(name, '%s') < 5", filter.StudyGroupName)
	}

	if filter.Location != "" {
		query += fmt.Sprintf(" AND levenshtein(location, '%s') < 5", filter.Location)
	}

	if filter.MeetingDate != "" {
		date := strings.Split(filter.MeetingDate, "T")[0]
		query += fmt.Sprintf(" AND to_char(meeting_date, 'YYYY-MM-DD') LIKE '%s'", date)
	}

	if filter.CourseCode != "" {
		query += fmt.Sprintf(" AND levenshtein(course ->> 'code', '%s') < 5", filter.CourseCode)
	}

	if filter.CourseName != "" {
		query += fmt.Sprintf(" AND levenshtein(course ->> 'name', '%s') < 5", filter.CourseName)
	}

	if filter.Instructor != "" {
		query += fmt.Sprintf(" AND levenshtein(course ->> 'instructor', '%s') < 5", filter.Instructor)
	}

	if filter.Term != "" {
		query += fmt.Sprintf(" AND levenshtein(course ->> 'term', '%s') < 5", filter.Term)
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, filter.PageIndex)

	if err := server.DB.Select(&studyGroups, query); err != nil {
		return studyGroups, http.StatusInternalServerError, errors.New(
			"unable to get study groups",
		)
	}

	return studyGroups, http.StatusOK, nil
}

func GetStudyGroupMembers(studyGroupID string) (interface{}, int, error) {
	var (
		studyGroup models.StudyGroup
		members []models.User
		waitlist []models.User
		users = map[string][]models.User{
			"members": []models.User{},
			"waitlist": []models.User{},
		}
	)
	errMsg := errors.New("unable to get study group members")

	getUsers := func (userIDsCSV, usersType string, sgUsers []models.User) error {
		if userIDsCSV == "" { return nil }

		query := "SELECT * FROM users WHERE id = "
		for i, userID := range strings.Split(userIDsCSV, ",") {
			if i == 0 {
				query += userID
				continue
			}

			query += " OR id = " + userID
		}

		if err := server.DB.Select(&sgUsers, query); err != nil {
			return errMsg
		}

		users[usersType] = sgUsers
		return nil
	}


	err := server.DB.Get(&studyGroup,
		`SELECT members, waitlist
		 FROM study_groups
		 WHERE id = $1`,
		studyGroupID,
	)

	switch {
		case err == sql.ErrNoRows:
			return nil, http.StatusNotFound, errors.New("study group doesn't exist")
		case err != nil:
			return nil, http.StatusInternalServerError, errMsg
	}

	if err = getUsers(studyGroup.Members.String, "members", members); err != nil {
		return users, http.StatusInternalServerError, err
	}

	if err = getUsers(studyGroup.Waitlist.String, "waitlist", waitlist); err != nil {
		return users, http.StatusInternalServerError, err
	}

	return users, http.StatusOK, nil
}

func CreateStudyGroup(studyGroup models.StudyGroup) (models.StudyGroup, int, error) {
	var newStudyGroup models.StudyGroup

	err := server.DB.Get(
	 &newStudyGroup,
	 `INSERT INTO study_groups
			(user_id, name, members_limit, available_spots, location, description, meeting_date, course, created_on, updated_on)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING *`,
			studyGroup.UserID,
			studyGroup.Name,
			studyGroup.MembersLimit,
			studyGroup.MembersLimit,
			studyGroup.Location,
			studyGroup.Description,
			studyGroup.MeetingDate,
			studyGroup.Course,
			time.Now(),
			time.Now(),
		)

		if err != nil {
			log.Println(err.Error())
			return newStudyGroup, http.StatusInternalServerError,
				errors.New("unable to create study group")
		}

		return newStudyGroup, http.StatusOK, nil
}

func UpdateStudyGroup(studyGroup models.StudyGroup) (models.StudyGroup, int, error) {
	var updatedStudyGroup models.StudyGroup

	_, err := server.DB.Exec(
	 `UPDATE study_groups
		SET
			name          = $1,
			members_limit = $2,
			description   = $3,
			meeting_date  = $4,
			location      = $5
		WHERE id = $6
		RETURNING *`,
		studyGroup.Name,
		studyGroup.MembersLimit,
		studyGroup.Description,
		studyGroup.MeetingDate,
		studyGroup.Location,
		studyGroup.ID,
	)

	if err != nil {
		return updatedStudyGroup, http.StatusInternalServerError,
			errors.New("unable to update study group")
	}

	return updatedStudyGroup, http.StatusOK, nil
}

func DeleteStudyGroup(studyGroupID, userID string) (int, error) {
	var studyGroup models.StudyGroup

	internalErr := func () (int, error) {
		return http.StatusInternalServerError, errors.New("unable delete study group")
	}

	err := server.DB.Get(
	 &studyGroup,
	 "SELECT members, waitlist FROM study_groups WHERE id = $1 AND user_id = $2",
		studyGroupID, userID,
	)

	if err != nil {	return internalErr() }

	var studyGroupUserIDs []string
	membersCSV := studyGroup.Members.String
	waitlistCSV := studyGroup.Waitlist.String

	if membersCSV != "" {
		studyGroupUserIDs = append(
			studyGroupUserIDs,
			strings.Split(membersCSV, ",")...
		)
	}

	if waitlistCSV != "" {
		studyGroupUserIDs = append(
			studyGroupUserIDs,
			strings.Split(waitlistCSV, ",")...
		)
	}

	if studyGroupUserIDs != nil {
		query := "SELECT id, study_groups, waitlists FROM users WHERE id = " + studyGroupUserIDs[0]

		for i := 1; i < len(studyGroupUserIDs); i++ {
			query += " OR id = " + studyGroupUserIDs[i]
		}

		var users []models.User
		if err = server.DB.Select(&users, query); err != nil {
			return internalErr()
		}

		{
			tx, err := server.DB.Begin()

			defer func() (int, error) {
				if err != nil {
					log.Println(err.Error())
					tx.Rollback()
					return internalErr()
				}

				return 0, nil
			}()

			for _, user := range users {
				user.LeaveStudyGroup(studyGroupID)

				_, err = tx.Exec(
				 "UPDATE users SET study_groups = $1, waitlists = $2 WHERE id = $3",
					user.StudyGroups,
					user.Waitlists,
					user.ID,
				)
			}

			_, err = tx.Exec("DELETE FROM study_groups WHERE id = $1", studyGroupID)
			err = tx.Commit()
		}
	}

	return http.StatusOK, nil
}

func JoinStudyGroup(studyGroupID, userID string) (models.StudyGroup, int, error) {
	var user models.User
	var studyGroup models.StudyGroup

	internalErr := func() (models.StudyGroup, int, error) {
		return studyGroup, http.StatusInternalServerError, errors.New("unable to join study group")
	}

	err := server.DB.Get(
		&studyGroup,
		"SELECT user_id, members, waitlist, available_spots FROM study_groups WHERE id = $1",
		studyGroupID,
	)

	switch {
	case err == sql.ErrNoRows:
		return studyGroup, http.StatusNotFound, errors.New("study group not found")
	case err != nil:
		return internalErr()
	}

	err = server.DB.Get(
		&user,
		"SELECT study_groups, waitlists FROM users WHERE id = $1",
		userID,
	)

	switch {
	case err == sql.ErrNoRows:
		return studyGroup, http.StatusNotFound, errors.New("user not found")
	case err != nil:
		return internalErr()
	}

	if err = studyGroup.AddUserToWaitlist(userID); err != nil {
		return studyGroup, http.StatusForbidden, err
	}
	if err = user.AddStudyGroupToWaitlists(studyGroupID); err != nil {
		return studyGroup, http.StatusForbidden, err
	}

	{
		tx, err := server.DB.Begin()
		if err != nil {
			return internalErr()
		}

		defer func() (models.StudyGroup, int, error) {
			if err != nil {
				log.Println(err.Error())
				tx.Rollback()
				return internalErr()
			}

			return studyGroup, 0, nil
		}()

		if studyGroup.Waitlist.String == "" {
			_, err = tx.Exec(
				"UPDATE study_groups SET waitlist = null, available_spots = $1 WHERE id = $2",
				studyGroup.AvailableSpots,
				studyGroupID,
			)
		} else {
			_, err = tx.Exec(
				"UPDATE study_groups SET waitlist = $1, available_spots = $2 WHERE id = $3",
				studyGroup.Waitlist.String,
				studyGroup.AvailableSpots,
				studyGroupID,
			)
		}

		if user.Waitlists.String == "" {
			_, err = tx.Exec(
			 "UPDATE users SET waitlists = null WHERE id = $1",
				userID,
			)
		} else {
			_, err = tx.Exec(
			 "UPDATE users SET waitlists = $1 WHERE id = $2",
				user.Waitlists.String,
				userID,
			)
		}

		err = tx.Commit()
	}

	return studyGroup, http.StatusOK, nil
}

func MoveUserFromWaitlistToMembers(studyGroupID, userID string) (models.StudyGroup, int, error) {
	var user models.User
	var studyGroup models.StudyGroup

	internalErr := func() (models.StudyGroup, int, error) {
		return studyGroup, http.StatusInternalServerError, errors.New("unable to move user into members")
	}

	err := server.DB.Get(&studyGroup, "SELECT * FROM study_groups WHERE id = $1", studyGroupID)
	if err != nil {	return internalErr() }

	if err := studyGroup.MoveUserFromWaitlistToMembers(userID); err != nil {
		return studyGroup, http.StatusForbidden, err
	}

	err = server.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {	return internalErr() }

	{
		tx, err := server.DB.Begin()
		if err != nil {	return internalErr() }

		defer func() (models.StudyGroup, int, error) {
			if err != nil {
				log.Println(err.Error())
				tx.Rollback()

				return internalErr()
			}

			return studyGroup, 0, nil
		}()

		_, err = tx.Exec(
			"UPDATE study_groups SET members = $1, waitlist = $2 WHERE id = $3",
			studyGroup.Members,
			studyGroup.Waitlist,
			studyGroup.ID,
		)

		_, err = tx.Exec(
			"UPDATE users SET study_groups = $1, waitlists = $2 WHERE id = $3",
			user.StudyGroups,
			user.Waitlists,
			user.ID,
		)

		err = tx.Commit()
	}

	return studyGroup, http.StatusOK, nil
}

func LeaveStudyGroup(studyGroupID, userID string) (int, error) {
	var user models.User
	var studyGroup models.StudyGroup

	internalErr := func() (int, error) {
		return http.StatusInternalServerError, errors.New("unable to leave study group")
	}

	err := server.DB.Get(
		&studyGroup,
		"SELECT user_id, members, waitlist, available_spots FROM study_groups WHERE id = $1",
		studyGroupID,
	)

	switch {
	case err == sql.ErrNoRows:
		return http.StatusNotFound, errors.New("study group not found")
	case err != nil:
		return internalErr()
	}

	err = server.DB.Get(
		&user,
		"SELECT study_groups, waitlists FROM users WHERE id = $1",
		userID,
	)

	switch {
	case err == sql.ErrNoRows:
		return http.StatusNotFound, errors.New("user not found")
	case err != nil:
		return internalErr()
	}

	sgColumnName, sgColumnVal, err := studyGroup.RemoveUser(userID)
	if err != nil {
		return http.StatusForbidden, err
	}

	uColumnName, uColumnVal, err := user.LeaveStudyGroup(studyGroupID)
	if err != nil {
		return http.StatusForbidden, err
	}

	{
		tx, err := server.DB.Begin()
		if err != nil {
			return internalErr()
		}

		defer func() (int, error) {
			if err != nil {
				log.Println(err.Error())
				tx.Rollback()
				return internalErr()
			}

			return 0, nil
		}()

		if sgColumnVal.String == "" {
			_, err = tx.Exec(
				"UPDATE study_groups SET "+sgColumnName+" = null, available_spots = $1 WHERE id = $2",
				studyGroup.AvailableSpots,
				studyGroupID,
			)
		} else {
			_, err = tx.Exec(
				"UPDATE study_groups SET "+sgColumnName+" = $1, available_spots = $2 WHERE id = $3",
				sgColumnVal.String,
				studyGroup.AvailableSpots,
				studyGroupID,
			)
		}

		if uColumnVal.String == "" {
			_, err = tx.Exec("UPDATE users SET "+uColumnName+" = null WHERE id = $1",
				userID,
			)
		} else {
			_, err = tx.Exec("UPDATE users SET "+uColumnName+" = $1 WHERE id = $2",
				uColumnVal.String,
				userID,
			)
		}

		err = tx.Commit()
	}

	return http.StatusOK, nil
}
