package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/quantum0hound/gochat/internal/models"
	"github.com/quantum0hound/gochat/internal/repository"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	accessTokenTTL  = 5 * time.Minute
	refreshTokenTTL = 24 * 30 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

type AuthService struct {
	userProvider    repository.UserProvider
	sessionProvider repository.SessionProvider
}

func NewAuthService(userProvider repository.UserProvider, sessionProvider repository.SessionProvider) *AuthService {
	return &AuthService{
		userProvider:    userProvider,
		sessionProvider: sessionProvider,
	}
}

func (s *AuthService) CreateUser(user *models.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.userProvider.Create(user)
}

func (s *AuthService) GenerateAccessToken(user *models.User) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			user.Id,
			user.Username,
		},
	)
	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

func (s *AuthService) GenerateAccessTokenId(userId int) (string, error) {

	user, err := s.userProvider.GetById(userId)
	if err != nil {
		return "", err
	}
	return s.GenerateAccessToken(user)
}

func (s *AuthService) CreateSession(userId int, fingerprint string) (*models.Session, error) {
	if len(fingerprint) == 0 {
		return nil, errors.New("empty fingerprint field")
	}
	session := models.Session{
		UserId:       userId,
		RefreshToken: uuid.New().String(),
		ExpiresIn:    time.Now().Add(refreshTokenTTL),
		Fingerprint:  fingerprint,
	}
	err := s.sessionProvider.Create(&session)
	return &session, err

}

func (s *AuthService) RefreshSession(refreshToken, fingerprint string) (*models.Session, error) {
	if len(fingerprint) == 0 {
		return nil, errors.New("empty fingerprint field")
	}
	session, err := s.sessionProvider.Get(refreshToken)
	if err != nil {
		return nil, errors.New("failed to get refresh session")
	}
	logrus.Debugf("expires=%d, now=%d", session.ExpiresIn.Unix(), time.Now().Unix())
	// if refresh token has expired, delete the session
	if session.ExpiresIn.Unix() < time.Now().Unix() {
		errMessage := "refresh token has expired"
		err = s.sessionProvider.Delete(refreshToken)
		if err != nil {
			errMessage += ", failed to delete session : " + err.Error()
		}
		return nil, errors.New(errMessage)
	} else if session.ExpiresIn.Unix() < time.Now().Add(10*time.Minute).Unix() { //session is about to expire, create new, delete old
		err = s.sessionProvider.Delete(refreshToken)
		if err != nil {
			return nil, errors.New("failed to delete refresh session: " + err.Error())
		}
		session, err = s.CreateSession(session.UserId, fingerprint)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}

	//if we received a new fingerprint, create a new session and use it
	//todo: handle multiple fingerprint usages
	if session.Fingerprint != fingerprint {
		session, err = s.CreateSession(session.UserId, fingerprint)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	}
	return session, err

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

func (s *AuthService) GetUser(username, password string) (*models.User, error) {
	return s.userProvider.Get(username, s.generatePasswordHash(password))
}
