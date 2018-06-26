package server

import (
  "errors"
  "os"
  "time"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/badoux/checkmail"
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
)


var (
  DB         *sqlx.DB
  S3Uploader *s3manager.Uploader
  S3Service  *s3.S3

  JWTSigningKey = []byte(os.Getenv("JWT_SIGNING_TOKEN"))
)

const (
  AWSRegion   = "us-west-1"
  S3Bucket    = "study-groups"
  S3BucketURL = "https://s3-us-west-1.amazonaws.com/study-groups/"
)

func InitServer() error {
  var err error

  DB, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
  if err != nil { return err }

  sess, err := session.NewSession(&aws.Config{
    Credentials: credentials.NewStaticCredentials(
      os.Getenv("AWS_AKID"),
      os.Getenv("AWS_SECRET"),
      "",
    ),
    Region: aws.String(AWSRegion),
  })
  if err != nil { return err }

  S3Uploader = s3manager.NewUploader(sess)

  sess, err = session.NewSession(&aws.Config{
    Credentials: credentials.NewStaticCredentials(
      os.Getenv("AWS_AKID"),
      os.Getenv("AWS_SECRET"),
      "",
    ),
    Region: aws.String(AWSRegion),
  })
  if err != nil { return err }

  S3Service = s3.New(sess)

  return nil
}

func GenerateAuthToken() (string, error) {
  token := jwt.New(jwt.SigningMethodHS256)

  token.Claims = jwt.MapClaims{
    "exp": time.Now().Add(time.Hour * 730).Unix(), // ~ 1 month
    "iat": time.Now().Unix(),
  }

  tokenString, err := token.SignedString(JWTSigningKey)

  if err != nil { return tokenString, errors.New("error while signing auth token") }

  return tokenString, nil
}

func ValidateEmail(email string) (errMsg error) {
  errMsg = errors.New("invalid email address")

  if err := checkmail.ValidateFormat(email); err != nil {
    return errMsg
  }

  if err := checkmail.ValidateHost(email); err != nil {
    return errMsg
  }

  err := checkmail.ValidateHost(email)
  if _, ok := err.(checkmail.SmtpError); ok && err != nil {
    return errMsg
  }

  return nil
}

func Respond(c *gin.Context, data interface{}, message string, httpStatus int) {
  var success bool

  // HTTP Status Code >= 400 indicates error
  if httpStatus < 400 {
    success = true
  } else {
    success = false
  }

  c.JSON(
    httpStatus,
    map[string]interface{} {
      "data":    data,
      "status":  httpStatus,
      "message": message,
      "success": success,
    },
  )
}
