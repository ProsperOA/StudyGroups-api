package models

type Course struct {
  Code       string `json:"code"`
  Name       string `json:"name"`
  Instructor string `json:"instructor"`
  Term       string `json:"term"`
}
