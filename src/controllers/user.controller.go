package controllers

import (
  "database/sql"
  "errors"
  "fmt"
  "log"
  "mime/multipart"
  "net/http"
  "strings"
  "os"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  mailchimp "github.com/beeker1121/mailchimp-go"
  "github.com/beeker1121/mailchimp-go/lists/members"
  "github.com/prosperoa/study-groups/src/models"
  "github.com/prosperoa/study-groups/src/server"
  "github.com/prosperoa/study-groups/src/utils"
  "golang.org/x/crypto/bcrypt"
  "gopkg.in/guregu/null.v3"
)

func GetUser(userID string) (models.User, int, error) {
  var user models.User

  err := server.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userID)

  switch {
    case err == sql.ErrNoRows:
      return user, http.StatusNotFound, errors.New("user not found")
    case err != nil:
      return user, http.StatusInternalServerError, errors.New("unable to get user")
  }

  return user, http.StatusOK, nil
}

func GetUsers(page, pageSize int) ([]models.User, int, error) {
  var users []models.User

  err := server.DB.Select(&users, "SELECT * FROM users LIMIT $1 OFFSET $2",
    pageSize, pageSize * page,
  )

  switch {
    case err == sql.ErrNoRows, len(users) == 0:
      return users, http.StatusNotFound, errors.New("no users found")
    case err != nil:
      return users, http.StatusInternalServerError, errors.New("unable to get users")
  }

  return users, http.StatusOK, nil
}

func DeleteUser(userID string) (int, error) {
  var userEmail string

  err := server.DB.Get(&userEmail, "DELETE FROM users WHERE id = $1 RETURNING email", userID)

  switch {
    case err == sql.ErrNoRows:
      return http.StatusBadRequest, errors.New("account doesn't exist")
    case err != nil:
      return http.StatusInternalServerError, errors.New("unable delete account")
  }

  // remove user from mailchimp list
  var mailchimpID string

  params := &members.GetParams{Status: members.StatusSubscribed}
  listMembers, err := members.Get("4d6392ba4d", params)

  if err != nil {
    for _, v := range listMembers.Members {
      if v.EmailAddress == userEmail {
        mailchimpID = v.ID
        break
      }
    }

    if err = mailchimp.SetKey(os.Getenv("MAILCHIMP_API_KEY")); err != nil {
      log.Println(err.Error())
    }

    if err = members.Delete("4d6392ba4d", mailchimpID); err != nil {
      log.Println(err.Error())
    }
  } else {
    log.Println(err.Error())
  }

  return http.StatusOK, nil
}

func GetUserStudyGroups(userID string, page, pageSize int) ([]models.StudyGroup, int, error) {
  var studyGroups []models.StudyGroup

  err := server.DB.Select(
    &studyGroups,
    "SELECT * FROM study_groups WHERE user_id = $1 LIMIT $2 OFFSET $3",
    userID,
    pageSize,
    pageSize * page,
  )

  switch {
    case err == sql.ErrNoRows, len(studyGroups) == 0:
      return studyGroups, http.StatusNotFound, errors.New("no users study groups found")
    case err != nil:
      return studyGroups, http.StatusInternalServerError, errors.New("unable to get user's study groups")
  }

  return studyGroups, http.StatusOK, nil
}

func UploadAvatar(userID, ext string, image multipart.File) (string, int, error) {
  var avatarURL null.String
  errMsg := errors.New("unable to upload avatar")

  err := server.DB.Get(&avatarURL, "SELECT avatar FROM users WHERE id = $1", userID)

  switch {
    case err == sql.ErrNoRows:
      return avatarURL.String, http.StatusNotFound, errors.New("user not found")
    case err != nil:
      return avatarURL.String, http.StatusInternalServerError, errMsg
  }

  // delete old avatar
  if avatarURL.String != "" && !strings.Contains(avatarURL.String, "stock-avatar") {
    _, err = server.S3Service.DeleteObject(&s3.DeleteObjectInput{
      Bucket: aws.String(server.S3Bucket),
      Key:    aws.String(strings.TrimPrefix(avatarURL.String, server.S3BucketURL)),
    })

    if err != nil {
      return avatarURL.String, http.StatusInternalServerError, errMsg
    }
  }

  newAvatarFilename := fmt.Sprintf("%s-%s", userID, utils.RandString(16) + ext)

  result, err := server.S3Uploader.Upload(&s3manager.UploadInput{
    Body:   image,
    Bucket: aws.String(server.S3Bucket),
    Key:    aws.String("images/user-avatars/" + newAvatarFilename),
    ACL:    aws.String("public-read"),
  })

  if err != nil {
    return avatarURL.String, http.StatusInternalServerError, errMsg
  }

  res, err := server.DB.Exec("UPDATE users SET avatar = $1 WHERE id = $2",
    result.Location,
    userID,
  )

  rowsAffected, _ := res.RowsAffected()

  switch {
    case rowsAffected == 0:
      return avatarURL.String, http.StatusInternalServerError, errMsg
    case err != nil:
      return avatarURL.String, http.StatusInternalServerError, errMsg
  }

  return result.Location, http.StatusOK, nil
}

func ChangePassword(userID, currentPassword, desiredPassword string) (models.User, int, error) {
  var user models.User
  var currentPasswordHash string
  errMsg := errors.New("unable to change password")

  err := server.DB.Get(&currentPasswordHash, "SELECT password FROM users WHERE id = $1",
    userID,
  )
  if err != nil { return user, http.StatusInternalServerError, errMsg }

  err = bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(currentPassword))
  if err != nil { return user, http.StatusBadRequest, errors.New("incorrect password") }

  newPasswordHash, err  := bcrypt.GenerateFromPassword([]byte(desiredPassword), 6)
  if err != nil { return user, http.StatusInternalServerError, errMsg }

  err = server.DB.Get(&user, "UPDATE users SET password = $1 WHERE id = $2 RETURNING *",
    newPasswordHash, userID,
  )
  if err != nil { return user, http.StatusInternalServerError, errMsg }

  return user, http.StatusOK, nil
}
