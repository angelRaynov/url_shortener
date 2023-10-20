package model

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	UID      string `json:"uid"`
	Salt     string `json:"salt"`
}

type UserRequest struct {
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,min=6,max=64"`
	Email    string `json:"email" validate:"required,email"`
}

type EditUserRequest struct {
	Username string `json:"username" validate:"omitempty,max=64"`
	Password string `json:"password" validate:"omitempty,min=6,max=64"`
	Email    string `json:"email" validate:"omitempty,email"`
}
