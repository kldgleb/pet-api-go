package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"test-api/pkg/entity"
	"test-api/pkg/repository"
	"time"
)

const (
	salt       = "12hdhqwejkdbka"
	signingKey = "dsadaskm@122312dsa"
	tokenTTL   = 24 * time.Hour * 30
)

type AuthService struct {
	repo repository.Authorization
}

type customClaims struct {
	jwt.StandardClaims
	UserId int
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entity.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	logrus.Print(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetJWTByCredentials(username, password string) (string, error) {
	user, err := s.repo.GetUserByCredentials(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	claims := &customClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
