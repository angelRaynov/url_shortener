package usecase

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"url_shortener/helper"
	"url_shortener/infrastructure/config"
	"url_shortener/internal/model"
)

type logger interface {
	Debugw(msg string, keysAndValues ...interface{})
}
type userRepo interface {
	GetUserByUsername(username string) (*model.User, error)
	StoreUser(registerRequest *model.User, salt, hashedPass string) error
	UserExist(email string) bool
	EditUser(u *model.User, salt, uid string) error
}

type userUseCase struct {
	l   logger
	ur  userRepo
	cfg *config.Application
}

func NewUserUseCase(cfg *config.Application, l logger, ur userRepo) *userUseCase {
	return &userUseCase{
		l:   l,
		ur:  ur,
		cfg: cfg,
	}
}

func (au *userUseCase) GenerateJWT(creds model.AuthRequest) (string, error) {
	user, err := au.ur.GetUserByUsername(creds.Username)
	if err != nil {
		au.l.Debugw("getting user password", "error", err, "username", creds.Username)
		return "", fmt.Errorf("getting user password:%w", err)
	}

	passWithSalt := fmt.Sprintf("%s:%s", user.Salt, creds.Password)
	// Compare the entered password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passWithSalt))
	if err != nil {
		au.l.Debugw("comparing hash and password", "error", err, "username", creds.Username)
		return "", fmt.Errorf("comparing passwords:%w", err)
	}

	// Define the expiration time
	duration := time.Duration(au.cfg.JWTExpiration)
	expirationTime := time.Now().Add(duration * time.Hour)

	// Create the custom claims
	claims := model.Claims{
		Username: user.Username,
		Email:    user.Email,
		UID:      user.UID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a new token with custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(au.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("signing token:%w", err)
	}

	return tokenString, nil
}

func (au *userUseCase) RegisterUser(u *model.User) error {
	salt, err := helper.GenerateSalt()
	if err != nil {
		au.l.Debugw("generating salt", "error", err)
		return fmt.Errorf("generating salt:%w", err)
	}

	if au.ur.UserExist(u.Email) {
		au.l.Debugw("user already exist", "email", u.Email)
		return fmt.Errorf("user with email %s already exist", u.Email)
	}

	// Concatenate the salt and the password
	saltedPassword := fmt.Sprintf("%s:%s", salt, u.Password)
	// Hash the salted password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		au.l.Debugw("hashing password", "error", err, "email", u.Email)
		return fmt.Errorf("hashing password:%w", err)
	}

	err = au.ur.StoreUser(u, salt, string(hashedPassword))
	if err != nil {
		au.l.Debugw("storing user", "error", err)
		return fmt.Errorf("storing user:%w", err)
	}
	return nil
}

func (au *userUseCase) EditUser(u *model.User, uid string) error {
	salt, err := helper.GenerateSalt()
	if err != nil {
		au.l.Debugw("generating salt", "error", err)
		return fmt.Errorf("generating salt:%w", err)
	}

	if u.Password != "" {
		// Concatenate the salt and the password
		saltedPassword := fmt.Sprintf("%s:%s", salt, u.Password)
		// Hash the salted password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
		if err != nil {
			au.l.Debugw("hashing password", "error", err)
			return fmt.Errorf("hashing password:%w", err)
		}
		u.Password = string(hashedPassword)
	}

	err = au.ur.EditUser(u, salt, uid)
	if err != nil {
		au.l.Debugw("storing user", "error", err)
		return fmt.Errorf("storing user:%w", err)
	}
	return nil
}
