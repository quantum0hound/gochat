package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/repository"
	"os"
	"time"
)

const (
	tokenTTL = 1000 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

type AuthService struct {
	userProvider repository.UserProvider
}

func NewAuthService(userProvider repository.UserProvider) *AuthService {
	return &AuthService{userProvider: userProvider}
}

func (s *AuthService) CreateUser(user *models.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.userProvider.Create(user)
}

func (s *AuthService) GenerateAccessToken(username, password string) (string, error) {
	user, err := s.userProvider.Get(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			user.Id,
			username,
		},
	)
	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

func (s *AuthService) GenerateRefreshToken() string {
	return uuid.New().String()
}

func (s *AuthService) ParseAccessToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(os.Getenv("SIGNING_KEY")), nil
		},
	)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type of *tokenClaims")
	}
	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}
