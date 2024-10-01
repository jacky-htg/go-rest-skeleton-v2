package handler

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"rest-skeleton/internal/dto"
	"rest-skeleton/internal/pkg/logger"
	"rest-skeleton/internal/usecase"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Auths struct {
	Log *logger.Logger
	DB  *sql.DB
}

// @Summary Login
// @Description Login to the system
// @host   localhost:8081
// @ID login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param Idempotency-Key header string true "Idempotency-Key"
// @Param login body dto.LoginRequest true "Login"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /login [post]
func (h *Auths) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "LoginHandler")
	defer span.End()

	switch ctx.Err() {
	case context.Canceled:
		h.Log.Error(ctx, context.Canceled)
		http.Error(w, "Request is canceled", http.StatusExpectationFailed)
		return
	case context.DeadlineExceeded:
		h.Log.Error(ctx, context.DeadlineExceeded)
		http.Error(w, "Deadline is exceeded", http.StatusExpectationFailed)
		return
	default:
	}

	var loginRequest dto.LoginRequest

	defer r.Body.Close()
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("email", loginRequest.Email))
	if err := loginRequest.Validate(); err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var authUC = usecase.AuthUC{Log: h.Log, DB: h.DB}
	token, statusCode, err := authUC.Login(r.Context(), loginRequest)
	if err != nil {
		http.Error(w, "Login failed", statusCode)
		return
	}

	response := dto.LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := sonic.ConfigDefault.NewEncoder(w).Encode(response); err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
