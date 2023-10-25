package http

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"url_shortener/internal/model"
)

type logger interface {
	Debugw(msg string, keysAndValues ...interface{})
}
type jwtGenerator interface {
	GenerateJWT(creds model.AuthRequest) (string, error)
}

type regEditor interface {
	RegisterUser(u *model.User) error
	EditUser(u *model.User, uid string) error
}
type userUseCase interface {
	jwtGenerator
	regEditor
}

type userHandler struct {
	l           logger
	userUseCase userUseCase
}

func NewAuthHandler(l logger, userUseCase userUseCase) *userHandler {
	return &userHandler{
		l:           l,
		userUseCase: userUseCase,
	}
}

func (uh *userHandler) Authenticate(c *gin.Context) {
	var creds model.AuthRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		uh.l.Debugw("binding user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tokenString, err := uh.userUseCase.GenerateJWT(creds)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			uh.l.Debugw("wrong credentials", "error", err, "username", creds.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			uh.l.Debugw("wrong credentials", "error", err, "username", creds.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		uh.l.Debugw("generating jwt", "error", err, "username", creds.Username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Return the token in the response
	uh.l.Debugw("generated jwt", "username", creds.Username)
	c.JSON(http.StatusOK, model.AuthResponse{Token: tokenString})
}

func (uh *userHandler) Register(c *gin.Context) {
	var rr model.UserRequest
	if err := c.ShouldBindJSON(&rr); err != nil {
		uh.l.Debugw("binding user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	v := validator.New()

	err := v.Struct(rr)
	if err != nil {
		uh.l.Debugw("validating user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	um := mapUserToModel(rr)

	err = uh.userUseCase.RegisterUser(um)
	if err != nil {
		uh.l.Debugw("registering user:", "username", rr.Username, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	uh.l.Debugw("user registered successfully", "username", rr.Username)
	c.JSON(http.StatusCreated, nil)
}

func (uh *userHandler) Edit(c *gin.Context) {
	var er model.EditUserRequest
	if err := c.ShouldBindJSON(&er); err != nil {
		uh.l.Debugw("binding user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	v := validator.New()

	err := v.Struct(er)
	if err != nil {
		uh.l.Debugw("validating edit request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isEditRequestEmpty(er) {
		uh.l.Debugw("updating user: empty edit request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user := mapEditRequestToModel(er)

	uid, ok := c.Get("user_uid")
	if !ok {
		uh.l.Debugw("updating user: missing uid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID, isString := uid.(string)
	if !isString {
		uh.l.Debugw("updating user: uid not a string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = uh.userUseCase.EditUser(user, userID)
	if err != nil {
		uh.l.Debugw("updating user:", "username", er.Username, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update user failed"})
		return
	}

	uh.l.Debugw("user updated successfully", "username", er.Username)
	c.JSON(http.StatusCreated, nil)
}

func mapUserToModel(ur model.UserRequest) *model.User {
	return &model.User{
		Username: ur.Username,
		Password: ur.Password,
		Email:    ur.Email,
	}
}

func mapEditRequestToModel(ur model.EditUserRequest) *model.User {
	return &model.User{
		Username: ur.Username,
		Password: ur.Password,
		Email:    ur.Email,
	}
}
func isEditRequestEmpty(er model.EditUserRequest) bool {
	if er.Username == "" && er.Email == "" && er.Password == "" {
		return true
	}
	return false
}
