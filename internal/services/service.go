package services

const DefaultPageSize = 20


type Pagination struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
