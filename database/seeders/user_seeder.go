package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunUserSeeder(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := models.User{
		UUID:        uuid.New(),
		Name:        "Admin",
		Username:    "admin",
		Password:    string(password),
		PhoneNumber: "081234567890",
		Email:       "admin@gmail.com",
		RoleId:      constants.Admin,
	}

	err := db.FirstOrCreate(&users, models.User{Email: users.Email}).Error
	if err != nil {
		logrus.Errorf("failed to seed user: %v", err)
		panic(err)
	}

	logrus.Infof("success seeding user %v", users.Username)
}
