package dto

import "errors"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() error {
	if len(l.Email) == 0 {
		return errors.New("email is required")
	}

	if len(l.Password) == 0 {
		return errors.New("password is required")
	}

	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}
