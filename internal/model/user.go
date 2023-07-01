package model

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UID      string `json:"uid"`
}
