package dto

import (
	"errors"
	"regexp"
	"rest-skeleton/internal/model"
)

type UserCreateRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

func (u *UserCreateRequest) Validate() error {
	if len(u.Name) == 0 {
		return errors.New("name is required")
	}

	if len(u.Email) == 0 {
		return errors.New("email is required")
	}

	if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, u.Email); !match {
		return errors.New("email harus valid")
	}

	if len(u.Password) == 0 {
		return errors.New("password is required")
	}

	if len(u.Password) < 10 {
		return errors.New("password minimal 10 character")
	}

	if match, _ := regexp.MatchString(`[a-z]`, u.Password); !match {
		return errors.New("password harus mengandung 1 huruf kecil")
	}

	if match, _ := regexp.MatchString(`[A-Z]`, u.Password); !match {
		return errors.New("password harus mengandung 1 huruf besar")
	}

	if match, _ := regexp.MatchString(`[0-9]`, u.Password); !match {
		return errors.New("password harus mengandung 1 angka")
	}

	if match, _ := regexp.MatchString(`[^a-zA-Z0-9]`, u.Password); !match {
		return errors.New("password harus mengandung 1 karakter khusus")
	}

	if len(u.RePassword) == 0 {
		return errors.New("re_password is required")
	}

	if u.Password != u.RePassword {
		return errors.New("password and re_password not match")
	}

	return nil
}

func (u *UserCreateRequest) ToEntity() model.User {
	return model.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserUpdateRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *UserUpdateRequest) Validate(id int64) error {
	if id != u.ID {
		return errors.New("id not match with user id")
	}

	if len(u.Name) == 0 {
		return errors.New("name is required")
	}

	return nil
}

func (u *UserUpdateRequest) ToEntity() model.User {
	return model.User{
		ID:   u.ID,
		Name: u.Name,
	}
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *UserResponse) FromEntity(user model.User) {
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
}

func (u *UserResponse) ListFromEntity(users []model.User) []UserResponse {
	var list []UserResponse = make([]UserResponse, 0)
	for _, user := range users {
		var userResponse UserResponse
		userResponse.FromEntity(user)
		list = append(list, userResponse)
	}
	return list
}
