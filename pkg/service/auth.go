package service

import (
	"crypto/sha1"
	"documentStorage/models"
	"documentStorage/pkg"
	"documentStorage/pkg/repository"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	salt       = "anjn4f3krnvk35t9enkrjnv48jnc849"
	signingKey = "nvjsdni489jvdu3h8udjn49ghhnfo4jg"
	tokenTTl   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (string, error) {

	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	val, err := s.repo.GetToken(accessToken)
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, pkg.NewErrorResponse(500, "error accessing redis")
	} else if val == "blacklisted" {
		return 0, pkg.NewErrorResponse(401, "token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkg.NewErrorResponse(401, "invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, pkg.NewErrorResponse(500, "token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) Logout(token string) error {
	return s.repo.CreateToken(token)
}
