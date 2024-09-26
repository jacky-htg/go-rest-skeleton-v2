package handler

import (
	"context"
	"net/http"
	"rest-skeleton/dto"
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/usecase"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
)

type Auths struct {
	Log *logger.Logger
	DB  *database.Database
}

func (h *Auths) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch r.Context().Err() {
	case context.Canceled:
		h.Log.Error.Println("Request is canceled")
		http.Error(w, "Request is canceled", http.StatusExpectationFailed)
		return
	case context.DeadlineExceeded:
		h.Log.Error.Println("deadline is exceeded")
		http.Error(w, "Deadline is exceeded", http.StatusExpectationFailed)
		return
	default:
	}

	var loginRequest dto.LoginRequest

	defer r.Body.Close()
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := loginRequest.Validate(); err != nil {
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
		h.Log.Error.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
