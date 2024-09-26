package handler

import (
	"context"
	"fmt"
	"net/http"
	"rest-skeleton/dto"
	"rest-skeleton/model"
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/logger"
	"rest-skeleton/pkg/redis"
	"rest-skeleton/repository"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Log   *logger.Logger
	DB    *database.Database
	Cache *redis.Cache
}

func (h *Users) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB.Conn}
	users, err := userRepo.List(r.Context(), ps.ByName("search"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var usersResponse dto.UserResponse
	response := usersResponse.ListFromEntity(users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := sonic.ConfigDefault.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Users) GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "please supply a valid id", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("users.%d", id)
	if cacheValue, isExist := h.Cache.Get(r.Context(), key); isExist {
		response := cacheValue.(dto.UserResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := sonic.ConfigDefault.NewEncoder(w).Encode(response); err != nil {
			h.Log.Error.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB.Conn}
	userRepo.UserEntity = model.User{ID: int64(id)}
	err = userRepo.Find(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var response dto.UserResponse
	response.FromEntity(userRepo.UserEntity)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := sonic.ConfigDefault.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.Cache.Add(r.Context(), key, response)
}

func (h *Users) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	var userRequest dto.UserCreateRequest
	defer r.Body.Close()
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := userRequest.Validate(); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB.Conn}
	userRepo.UserEntity = userRequest.ToEntity()
	password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userRepo.UserEntity.Password = string(password)

	if err := userRepo.Save(r.Context()); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var response dto.UserResponse
	response.FromEntity(userRepo.UserEntity)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := sonic.ConfigDefault.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Users) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	idstr := ps.ByName("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "please supply a valid id", http.StatusBadRequest)
		return
	}

	var userRequest dto.UserUpdateRequest
	defer r.Body.Close()
	err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := userRequest.Validate(int64(id)); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB.Conn}
	userRepo.UserEntity = userRequest.ToEntity()
	if err := userRepo.Update(r.Context()); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var response dto.UserResponse
	response.FromEntity(userRepo.UserEntity)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := sonic.ConfigDefault.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Users) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	idstr := ps.ByName("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "please supply a valid id", http.StatusBadRequest)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB.Conn}
	userRepo.UserEntity = model.User{ID: int64(id)}
	if err := userRepo.Delete(r.Context()); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
