package rest

import (
	"context"
	"log"
	"net/http"
	"repoboost/internal/httputil"
	"repoboost/internal/user/model"
	"repoboost/internal/user/service"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserHandler struct {
	svc service.Service
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		svc: service.New(db),
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.svc.GetUser(context.TODO(), uint(id))
	if err != nil {
		if err := httputil.JSONError(w, r, http.StatusInternalServerError, err); err != nil {
			log.Println(err)
		}
	}
	if err := httputil.JSON(w, r, http.StatusOK, user); err != nil {
		log.Println(err)
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetUsers(context.TODO())
	if err != nil {
		if err := httputil.JSON(w, r, http.StatusInternalServerError, map[string]string{"error": err.Error()}); err != nil {
			log.Println(err)
		}
	}
	if err := httputil.JSON(w, r, http.StatusOK, users); err != nil {
		log.Println(err)
	}
}

type userRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	BranchID uint   `json:"branchId"`
	Role     string `json:"role"`
}

func (u userRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(5, 60)),
		validation.Field(&u.Username, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 20)),
	)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var u userRequest
	if err := httputil.Bind(w, r, &u); err != nil {
		if err := httputil.JSONError(w, r, http.StatusBadRequest, err); err != nil {
			log.Println(err)
		}
	}
	if err := validation.Validate(u); err != nil {
		if err := httputil.JSONError(w, r, http.StatusBadRequest, err); err != nil {
			log.Println(err)
		}
	}

	user := model.User{
		Name:     u.Name,
		Username: u.Username,
		Password: u.Password,
	}

	if err := h.svc.CreateUser(context.TODO(), &user); err != nil {
		if err == service.ErrInvalidRequest {
			if err := httputil.JSONMessage(w, r, http.StatusBadRequest, "invalid credentials"); err != nil {
				log.Println(err)
			}
		}

		if err := httputil.JSONError(w, r, http.StatusInternalServerError, err); err != nil {
			log.Println(err)
		}
	}

	if err := httputil.JSON(w, r, http.StatusOK, nil); err != nil {
		log.Println(err)
	}
}
