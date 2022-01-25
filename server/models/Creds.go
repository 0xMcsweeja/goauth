package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID          uint32      `gorm:"primary_key" json:"id"`
	Nickname    string      `gorm:"size:255;not null" json:"nickname"`
	Email       string      `gorm:"size:100;not null" json:"email"`
	Password    string      `gorm:"size:100;not null;" json:"password"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Credentials Credentials `gorm:"foreignKey:Secret;references:ID" json:"credentials"`
}

type Credentials struct {
	Credential_type string `json:"credential_type,omitempty"`
	Secret          string `json:"secret,omitempty"`
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}
