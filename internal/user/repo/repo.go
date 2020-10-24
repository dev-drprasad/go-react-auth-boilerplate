package repo

import (
	"context"
	"database/sql"
	"repoboost/internal/user/model"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var ErrNoUserFound = errors.New("no user found")
var ErrFKViolation = errors.New("violates foreign key constraint")

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return Repository{db: db}
}

func (repo Repository) GetUsers(ctx context.Context) ([]*model.User, error) {

	var users []*model.User

	stmt := `SELECT id, name, username FROM users`
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Username); err != nil {
			break
		}
		users = append(users, &user)
	}
	if err != nil {
		return users, errors.Wrap(err, "Failed to scan rows")
	}

	return users, nil
}

func (repo Repository) GetUser(ctx context.Context, id uint) (*model.User, error) {

	stmt := `SELECT id, name, username FROM users WHERE id = $1`

	var user model.User
	err := repo.db.QueryRow(ctx, stmt, id).Scan(&user.ID, &user.Name, &user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &user, nil
}

func (repo Repository) CreateUser(ctx context.Context, u *model.User) error {
	stmt := `
		INSERT INTO users
			(name, username, password)
		VALUES
			($1, $2, $3, crypt($4, gen_salt('bf')), $5)
		RETURNING id`
	var userID uint64

	err := repo.db.QueryRow(ctx, stmt, u.Name, u.Username, u.Password).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), ErrFKViolation.Error()) {
			return ErrFKViolation
		}
		return errors.Wrapf(err, "Failed to execute the query username=%v", u.Username)
	}

	return nil
}

func (repo Repository) GetUserByUsernamePassword(ctx context.Context, username string, password string) (*model.User, error) {
	stmt := `
		SELECT users.id, users.name, users.username
		FROM users
		WHERE username=$1 AND password = crypt($2, password)
	`

	var user model.User
	err := repo.db.QueryRow(ctx, stmt, username, password).Scan(&user.ID, &user.Name, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoUserFound
		}
		return nil, errors.Wrap(err, "Failed to execute the query")
	}

	return &user, nil
}
