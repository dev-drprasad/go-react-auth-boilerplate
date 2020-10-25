package service

import (
	"context"
	"fmt"
	"os"
	"repoboost/internal/user/model"
	"repoboost/internal/user/repo"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

var sugar string

func init() {
	var sugar string
	sugar = os.Getenv("JWT_SECRET")

	if sugar == "" {
		panic("env JWT_SECRET is empty and it is required by auth serice")
	}
}

type Service struct {
	repo repo.Repository
}

func New(db *pgxpool.Pool) Service {
	return Service{
		repo: repo.New(db),
	}
}

func generateUserToken(u *model.User, sugar string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = u.ID
	claims["name"] = u.Name
	claims["username"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(sugar))
}

func AuthCheck(tokenStr string) (*model.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, msg
		}
		return []byte(sugar), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error parsing token")

	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	user := model.User{
		ID:       uint(claims["id"].(float64)),
		Username: claims["username"].(string),
	}
	return &user, nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsernamePassword(ctx, username, password)
	if err != nil {
		if err == repo.ErrNoUserFound {
			return "", ErrInvalidCredentials
		}
		return "", err

	}
	return generateUserToken(user, sugar)
}
