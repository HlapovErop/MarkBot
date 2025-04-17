package models

import (
	"github.com/HlapovErop/MarkBot/src/internal/utils/toggles"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"slices"
)

// Константы для ролей
const (
	RoleStudent int64 = 1
	RoleTeacher int64 = 2
)

type InterfaceUser interface {
	GetRoles() []int64
	GetID() uint
}

type User struct {
	gorm.Model
	Name     string        `json:"name" gorm:"type:varchar(255);not null;default:null;index:idx_name_surname,unique"`
	Surname  string        `json:"surname" gorm:"type:varchar(255);not null;default:null;index:idx_name_surname,unique"`
	Email    string        `json:"email" gorm:"type:varchar(255);not null;default:null;unique_index"`
	Password string        `json:"-" gorm:"type:varchar(255);not null;default:null"`
	Points   int64         `json:"points" gorm:"type:int;not null;default:100"`
	Roles    pq.Int64Array `json:"roles" gorm:"type:integer[];not null;"`
}

func (u *User) GetRoles() []int64 {
	return u.Roles
}

func (u *User) GetID() uint {
	return u.ID
}

func IsTeacher(user InterfaceUser) bool {
	return slices.Contains(user.GetRoles(), RoleTeacher)
}

func IsStudent(user InterfaceUser) bool {
	return slices.Contains(user.GetRoles(), RoleStudent)
}

func (u *User) ValidateNameSurname() bool {
	allowedNames, ok := toggles.GetTogglesStorage().GetStringSlice("AllowedNames")
	if !ok {
		return true
	}

	fullName := u.Name + " " + u.Surname
	return slices.Contains(allowedNames, fullName)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
