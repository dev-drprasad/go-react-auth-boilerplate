package service

import (
	"context"
	"fmt"
	"log"

	"repoboost/internal/user/model"
	"repoboost/internal/user/repo"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRequest = errors.New("invalid request")

type Service struct {
	repo repo.Repository
}

func New(db *pgxpool.Pool) Service {
	return Service{
		repo: repo.New(db),
	}
}

func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) {
	log.Println("Reading user from repo. id=%v", id)
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) GetUsers(ctx context.Context) (users []*model.User, err error) {
	return s.repo.GetUsers(ctx)
}

func (s *Service) GetUserByUsernamePassword(ctx context.Context, username string, password string) (*model.User, error) {
	return s.repo.GetUserByUsernamePassword(ctx, username, password)
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if err == repo.ErrFKViolation {
			return ErrInvalidRequest
		}
		return fmt.Errorf("Failed to create user. %v", err)
	}
	return nil
}
