package response

import (
	"time"
)

type AdminData struct {
	IdMsb              uint64    `structs:"id_msb"`
	IdLsb              uint64    `structs:"id_lsb"`
	FirstName          string    `structs:"first_name"`
	LastName           string    `structs:"last_name"`
	Email              string    `structs:"email"`
	Password           string    `structs:"password"`
	Status             string    `structs:"status"`
	SecretKey          string    `structs:"secret_key"`
	CreatedAt          time.Time `structs:"created_at"`
	UpdatedAt          time.Time `structs:"updated_at"`
	ResetPasswordToken string    `structs:"reset_password_token"`
}

type AdminResponse struct {
	IdMsb     uint64    `structs:"id_msb"`
	IdLsb     uint64    `structs:"id_lsb"`
	FirstName string    `structs:"first_name"`
	LastName  string    `structs:"last_name"`
	Email     string    `structs:"email"`
	Status    string    `structs:"status"`
	CreatedAt time.Time `structs:"created_at"`
}

type AdminListResponse struct {
	IdMsb     uint64    `structs:"id_msb"`
	IdLsb     uint64    `structs:"id_lsb"`
	FirstName string    `structs:"first_name"`
	LastName  string    `structs:"last_name"`
	Email     string    `structs:"email"`
	Status    string    `structs:"status"`
	LastLogin *int      `structs:"last_login"`
	CreatedAt time.Time `structs:"created_at"`
}

type UserUUID struct {
	IdMsb uint64
	IdLsb uint64
}

type Login struct {
	Id             int64      `json:"id"`
	Name           string     ` json:"name"`
	Email          string     `json:"email"`
	Password       string     ` json:"password,omitempty"`
	VerifiedStatus int        `json:"verified_status"`
	Status         int        `json:"status" `
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type LoginResponse struct {
	Id        int64      `json:"id" sql:"index"`
	Name      string     `json:"name"`
	Email     string     ` json:"email"`
	Password  string     ` json:"password"`
	AdminType int64      ` json:"admin_type"`
	Status    int        `json:"status" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
