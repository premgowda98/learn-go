package models

type UserRequest struct {
	ID   int64 `json:"id"`
	Name string `json:"name" validate:"required"`
}

type Request struct {
	UserRequest
}
