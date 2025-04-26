package usersEntity

import (
	"go-auth/app/models"

	"gorm.io/gorm"
)

func GetFirstByEmail(email string, db *gorm.DB) models.Users{
	var users models.Users
	db.Raw("SELECT *from users where email=?", email).Scan(&users)
	return users
}


 