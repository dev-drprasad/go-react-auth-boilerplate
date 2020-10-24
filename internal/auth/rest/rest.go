package rest

import (
	"context"
	"log"
	"net/http"
	"repoboost/internal/auth/service"
	"repoboost/internal/httputil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthHandler struct {
	svc service.Service
}

func New(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		svc: service.New(db),
	}
}

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p loginPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Username, validation.Required),
		validation.Field(&p.Password, validation.Required),
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var payload loginPayload
	if err := httputil.Bind(w, r, &payload); err != nil {
		log.Println(err)
	}

	token, err := h.svc.Login(context.TODO(), payload.Username, payload.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			if err := httputil.JSONMessage(w, r, http.StatusUnauthorized, "invalid credentials"); err != nil {
				log.Println(err)
			}
			return
		}
		if err := httputil.JSONError(w, r, http.StatusInternalServerError, err); err != nil {
			log.Println(err)
		}
		return
	}
	if err := httputil.JSON(w, r, http.StatusOK, map[string]string{"token": token}); err != nil {
		log.Println(err)
		return
	}
}
