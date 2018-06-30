package emails

import (
  "bytes"
  "html/template"
  "errors"
  "net/smtp"
  "os"

  "github.com/jordan-wright/email"
)

var (
  newUserTpl = template.Must(template.New("new-user.html").ParseFiles(
    "./email-notifications/templates/new-user.html",
  ))
)

var errMsg = errors.New("unable to send email notification")

type emailUser struct {
  Name  string
  Email string
}

func NewUserNotification(userName, recipientEmail string) (errMsg error) {
  var buf bytes.Buffer

  data := emailUser{
    Name: userName,
    Email: recipientEmail,
  }

  if err := newUserTpl.Execute(&buf, &data); err != nil {
    return errMsg
  }

  email := &email.Email{
    To: []string{recipientEmail},
    From: "StudyGroups <studygroups.io@gmail.com>",
    Subject: "Welcome to Study Groups",
    HTML: []byte(buf.String()),
  }

  err := email.Send(
    os.Getenv("SMTP_HOST") + os.Getenv("SMTP_PORT"),
    smtp.PlainAuth(
      "",
      os.Getenv("SMTP_EMAIL"),
      os.Getenv("SMTP_PASSWORD"),
      os.Getenv("SMTP_HOST"),
    ),
  )

  if err != nil {
    return errMsg
  }

  return nil
}
