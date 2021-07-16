package model

import "time"

type Admin struct {
	Model
	Name             string    `gorm:"type:varchar(50); not null;" json:"name" validate:"required"`
	Email            string    `gorm:"type:varchar(50); not null;" json:"email" validate:"required,email"`
	Password         string    `gorm:"type:varchar(100); not null;" json:"password" validate:"required"`
	ResetToken       string    `gorm:"type:varchar(100); not null;" json:"reset_token"`
	ResetTokenExpiry time.Time `gorm:"type:datetime; DEFAULT:null" json:"reset_token_expiry"`
	VerifiedStatus   int       `gorm:"type:int(2); DEFAULT:2; comment:'1=>Verified, 2=> Unverified'"  json:"verified_status" sql:"index"`
	Status           int       `gorm:"type:int(2); DEFAULT:2; comment:'1=>Active, 2=> Inactive'"  json:"status" sql:"index"`
}

func (a *Admin) TableName() string {
	return "admin"
}
