package models

type BaseFilter struct {
	PageIndex int    `json:"page_index" validate:"min=0"`
	PageSize  int    `json:"page_size"  validate:"min=10,max=30"`
}