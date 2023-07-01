package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"url_shortener/internal/model"
)

func GetUserFromContext(c *gin.Context) (*model.User, error) {
	ownerID, exist := c.Get("owner_uid")
	if !exist {
		return &model.User{}, fmt.Errorf("uid is missing from context")
	}

	uid, isString := ownerID.(string)
	if !isString {
		return &model.User{}, fmt.Errorf("uid is not a string")

	}

	e, exist := c.Get("email")
	if !exist {
		return &model.User{}, fmt.Errorf("email is missing from context")
	}

	email, isString := e.(string)
	if !isString {
		return &model.User{}, fmt.Errorf("email is not a string")

	}

	u, exist := c.Get("username")
	if !exist {
		return &model.User{}, fmt.Errorf("username is missing from context")
	}

	username, isString := u.(string)
	if !isString {
		return &model.User{}, fmt.Errorf("username is not a string")

	}

	return &model.User{
		Username: username,
		Email:    email,
		UID:      uid,
	}, nil
}
