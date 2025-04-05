package database

import "github.com/HlapovErop/MarkBot/internal/models"

func seeds() {
	users := []models.User{
		{
			Name:     "epkhlapov",
			Email:    "xlapov21@mail.ru",
			Password: "123456789",
			Roles:    []int64{1},
		},
	}
}
