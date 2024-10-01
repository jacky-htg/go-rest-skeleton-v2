package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"rest-skeleton/internal/dto"
	"rest-skeleton/internal/model"
	"rest-skeleton/internal/pkg/httpresponse"
	"rest-skeleton/internal/pkg/logger"
	"rest-skeleton/internal/pkg/redis"
	"rest-skeleton/internal/repository"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

// Users handler
type Users struct {
	Log   *logger.Logger
	DB    *sql.DB
	Cache *redis.Cache
}

// @Security Bearer
// @Summary List Users
// @Description List Users
// @Tags Users
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UserResponse
// @Router /users [get]
func (h *Users) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "listUserHandler")
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

	span.SetAttributes(attribute.String("search", ps.ByName("search")))
	var httpres = httpresponse.Response{Cache: h.Cache}
	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB}
	users, err := userRepo.List(ctx, ps.ByName("search"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var usersResponse dto.UserResponse
	response := usersResponse.ListFromEntity(users)
	httpres.SetMarshal(ctx, w, http.StatusOK, response, "")
}

// @Security Bearer
// @Summary Get User By ID
// @Description Get User By ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [get]
func (h *Users) GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "GetUserByIdHandler")
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

	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr)
	span.SetAttributes(attribute.Int("id", id))
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "please supply a valid id", http.StatusBadRequest)
		return
	}

	httpres := httpresponse.Response{Cache: h.Cache}
	key := fmt.Sprintf("users.%d", id)
	if cacheValue, isExist := h.Cache.Get(ctx, key); isExist {
		span.SetAttributes(attribute.String("cache-key", key))
		httpres.Set(w, http.StatusOK, cacheValue)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB}
	userRepo.UserEntity = model.User{ID: int64(id)}
	err = userRepo.Find(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var response dto.UserResponse
	response.FromEntity(userRepo.UserEntity)
	httpres.SetMarshal(ctx, w, http.StatusOK, response, key)
}

// @Security Bearer
// @Summary Create User
// @Description Create User
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body dto.UserCreateRequest true "User to add"
// @Param Idempotency-Key header string true "Idempotency-Key"
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} dto.UserResponse
// @Router /users [post]
func (h *Users) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx = r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "CreateUserHandler")
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

	var httpres = httpresponse.Response{Cache: h.Cache}
	var userRequest dto.UserCreateRequest
	defer r.Body.Close()
	err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := userRequest.Validate(); err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB}
	userRepo.UserEntity = userRequest.ToEntity()
	password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userRepo.UserEntity.Password = string(password)

	if err := userRepo.Save(ctx); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var response dto.UserResponse
	response.FromEntity(userRepo.UserEntity)
	httpres.SetMarshal(ctx, w, http.StatusCreated, response, "")
}

// @Security Bearer
// @Summary Update User
// @Description Update User
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body dto.UserUpdateRequest true "User to update"
// @Param Idempotency-Key header string true "Idempotency-Key"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [put]
func (h *Users) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "UpdateUserHandler")
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

	idstr := ps.ByName("id")
	id, err := strconv.Atoi(idstr)
	span.SetAttributes(attribute.Int("id", id))
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "please supply a valid id", http.StatusBadRequest)
		return
	}

	var httpres = httpresponse.Response{Cache: h.Cache}
	var userRequest dto.UserUpdateRequest
	defer r.Body.Close()
	err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := userRequest.Validate(int64(id)); err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB}
	userRepo.UserEntity = userRequest.ToEntity()
	if err := userRepo.Update(ctx); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var response dto.UserResponse
	response.FromEntity(userRepo.UserEntity)
	httpres.SetMarshal(ctx, w, http.StatusOK, response, "")
	h.Cache.Del(ctx, fmt.Sprintf("users.%d", id))
}

// @Security Bearer
// @Summary Delete User By ID
// @Description Delete User By ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param Idempotency-Key header string true "Idempotency-Key"
// @Param Authorization header string true "Bearer token"
// @Success 204
// @Router /users/{id} [delete]
func (h *Users) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ctx = r.Context()
	ctx, span := otel.Tracer(os.Getenv("APP_NAME")).Start(ctx, "DeleteUserHandler")
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

	idstr := ps.ByName("id")
	id, err := strconv.Atoi(idstr)
	span.SetAttributes(attribute.Int("id", id))
	if err != nil {
		h.Log.Error(ctx, err)
		http.Error(w, "please supply a valid id", http.StatusBadRequest)
		return
	}

	var userRepo = repository.UserRepository{Log: h.Log, Db: h.DB}
	userRepo.UserEntity = model.User{ID: int64(id)}
	if err := userRepo.Delete(ctx); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.Cache.Del(ctx, fmt.Sprintf("users.%d", id))
}
