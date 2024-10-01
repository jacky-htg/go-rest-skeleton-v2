package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rest-skeleton/internal/dto"
	"rest-skeleton/internal/handler"
	"rest-skeleton/internal/pkg/myctx"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type UserCreateScenario struct {
	Name           string
	Data           dto.UserCreateRequest
	ExpectedErr    string
	StatusCode     int
	IdempotencyKey string
}

func getScenarios() []UserCreateScenario {
	return []UserCreateScenario{
		{
			Name: "Invalid Name",
			Data: dto.UserCreateRequest{
				Name:       "",
				Email:      "john.doe@example.com",
				Password:   "password123",
				RePassword: "password123",
			},
			ExpectedErr:    "Invalid input: name is required",
			StatusCode:     http.StatusBadRequest,
			IdempotencyKey: uuid.NewString(),
		},
		{
			Name: "Invalid Email",
			Data: dto.UserCreateRequest{
				Name:       "John Doe",
				Email:      "",
				Password:   "password123",
				RePassword: "password123",
			},
			ExpectedErr:    "Invalid input: email is required",
			StatusCode:     http.StatusBadRequest,
			IdempotencyKey: uuid.NewString(),
		},
		{
			Name: "Invalid Password",
			Data: dto.UserCreateRequest{
				Name:       "John Doe",
				Email:      "john.doe@example.com",
				Password:   "",
				RePassword: "password123",
			},
			ExpectedErr:    "Invalid input: password is required",
			StatusCode:     http.StatusBadRequest,
			IdempotencyKey: uuid.NewString(),
		},
		{
			Name: "Invalid Re-Password",
			Data: dto.UserCreateRequest{
				Name:       "John Doe",
				Email:      "john.doe@example.com",
				Password:   "Password123!",
				RePassword: "",
			},
			ExpectedErr:    "Invalid input: re_password is required",
			StatusCode:     http.StatusBadRequest,
			IdempotencyKey: uuid.NewString(),
		},
		{
			Name: "Invalid Name",
			Data: dto.UserCreateRequest{
				Name:       "",
				Email:      "john.doe@example.com",
				Password:   "password123",
				RePassword: "password123",
			},
			ExpectedErr:    "Invalid input: name is required",
			StatusCode:     http.StatusBadRequest,
			IdempotencyKey: uuid.NewString(),
		},
		{
			Name: "Invalid Weak Password",
			Data: dto.UserCreateRequest{
				Name:       "John Doe",
				Email:      "john.doe@example.com",
				Password:   "password123",
				RePassword: "password123",
			},
			ExpectedErr:    "Invalid input: password harus mengandung 1 huruf besar",
			StatusCode:     http.StatusBadRequest,
			IdempotencyKey: uuid.NewString(),
		},
		{
			Name: "Valid User",
			Data: dto.UserCreateRequest{
				Name:       "John Doe",
				Email:      "john.doe@example.com",
				Password:   "Password123!",
				RePassword: "Password123!",
			},
			ExpectedErr:    "",
			StatusCode:     http.StatusCreated,
			IdempotencyKey: uuid.NewString(),
		},
	}
}

func TestCreateUser(t *testing.T) {
	userHandler := handler.Users{DB: db, Log: log, Cache: cache}

	router := httprouter.New()
	router.POST("/users", mid.WrapMiddleware(publicMiddlewares, userHandler.Create))

	scenarios := getScenarios()
	var wg sync.WaitGroup
	for _, tt := range scenarios {
		wg.Add(1)
		go func(tt struct {
			Name           string
			Data           dto.UserCreateRequest
			ExpectedErr    string
			StatusCode     int
			IdempotencyKey string
		}) {
			defer wg.Done()
			t.Run(tt.Name, func(t *testing.T) {
				dataJSON, err := json.Marshal(tt.Data)
				if err != nil {
					t.Errorf("could not marshal user data: %v", err)
				}
				req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(dataJSON))
				if err != nil {
					t.Errorf("could not create request: %v", err)
				}
				ctx := context.WithValue(req.Context(), myctx.Key("user_id"), int64(425071490427828))
				//ctx = context.WithValue(ctx, myctx.Key("traceID"), "29c92c0f-b95e-458e-b23c-09eb5058e83b")

				req = req.WithContext(ctx)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Idempotency-Key", tt.IdempotencyKey)

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)
				status := rr.Code
				if status != tt.StatusCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tt.StatusCode)
				}
				if status != http.StatusCreated {
					if strings.TrimSpace(rr.Body.String()) != tt.ExpectedErr {
						t.Errorf("handler returned wrong error message: got %v want %v",
							rr.Body.String(), tt.ExpectedErr)
					}
				} else {
					var response map[string]interface{}
					if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
						t.Errorf("could not unmarshal response: %v", err)
					}
					if response["name"] != tt.Data.Name || response["email"] != tt.Data.Email {
						t.Errorf("handler returned wrong user: got %v want %v",
							response, tt.Data)
					}
				}
			})
		}(tt)
	}
	wg.Wait()
}
