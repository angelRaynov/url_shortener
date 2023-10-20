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

type ILogger interface {
	Debugw(msg string, keysAndValues ...interface{})
}
type jwtGenerator interface {
	GenerateJWT(creds model.AuthRequest) (string, error)
	RegisterUser(u *model.User) error
	EditUser(u *model.User, uid string) error
}

type authHandler struct {
	l           ILogger
	authUseCase jwtGenerator
}

func NewAuthHandler(l ILogger, authUsecase jwtGenerator) *authHandler {
	return &authHandler{
		l:           l,
		authUseCase: authUsecase,
	}
}

func (ah *authHandler) Authenticate(c *gin.Context) {
	var creds model.AuthRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		ah.l.Debugw("binding user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tokenString, err := ah.authUseCase.GenerateJWT(creds)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			ah.l.Debugw("wrong credentials", "error", err, "username", creds.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			ah.l.Debugw("wrong credentials", "error", err, "username", creds.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		ah.l.Debugw("generating jwt", "error", err, "username", creds.Username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Return the token in the response
	ah.l.Debugw("generated jwt", "username", creds.Username)
	c.JSON(http.StatusOK, model.AuthResponse{Token: tokenString})
}

func (ah *authHandler) Register(c *gin.Context) {
	var rr model.UserRequest
	if err := c.ShouldBindJSON(&rr); err != nil {
		ah.l.Debugw("binding user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	v := validator.New()

	err := v.Struct(rr)
	if err != nil {
		ah.l.Debugw("validating user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	um := mapUserToModel(rr)

	err = ah.authUseCase.RegisterUser(um)
	if err != nil {
		ah.l.Debugw("registering user:", "username", rr.Username, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	ah.l.Debugw("user registered successfully", "username", rr.Username)
	c.JSON(http.StatusCreated, nil)
}

func (ah *authHandler) Edit(c *gin.Context) {
	var er model.EditUserRequest
	if err := c.ShouldBindJSON(&er); err != nil {
		ah.l.Debugw("binding user request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	v := validator.New()

	err := v.Struct(er)
	if err != nil {
		ah.l.Debugw("validating edit request:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isEditRequestEmpty(er) {
		ah.l.Debugw("updating user: empty edit request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user := mapEditRequestToModel(er)

	uid, ok := c.Get("user_uid")
	if !ok {
		ah.l.Debugw("updating user: missing uid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID, isString := uid.(string)
	if !isString {
		ah.l.Debugw("updating user: uid not a string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = ah.authUseCase.EditUser(user, userID)
	if err != nil {
		ah.l.Debugw("updating user:", "username", er.Username, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update user failed"})
		return
	}

	ah.l.Debugw("user updated successfully", "username", er.Username)
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
