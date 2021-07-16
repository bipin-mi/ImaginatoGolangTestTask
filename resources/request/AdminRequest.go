package request

import (
	"time"
)

type TokenDataRequest struct {
	ID int64 `json:"id"`
}

type DeleteAdminRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type CreateAdminRequest struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name" validate:"required" conform:"name"`
	Email            string    `json:"email" validate:"required,email" conform:"email"`
	Password         string    `json:"password" validate:"required"`
	ResetToken       string    `json:"reset_token"`
	ResetTokenExpiry time.Time `json:"reset_token_expiry"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordData struct {
	Id               int64     `json:"id"`
	Email            string    `json:"email" validate:"required,email"`
	ResetToken       string    `json:"reset_token"`
	ResetTokenExpiry time.Time `json:"reset_token_expiry"`
}

type ResetPasswordRequest struct {
	ResetToken string `json:"reset_token" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type ResetPasswordData struct {
	Id                 uint64    `json:"id"`
	ResetPasswordToken string    `json:"reset_token"`
	Password           string    `json:"password"`
	ResetTokenExpiry   time.Time `json:"reset_token_expiry"`
}

type VerifyEmailRequest struct {
	ResetToken string `json:"reset_token" validate:"required"`
}
