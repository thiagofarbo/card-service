package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jwt-authentication/helpers"
	"jwt-authentication/models"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := "root:admin@tcp(127.0.0.1:3306)/card?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed database connection")
	}

	DB = database
}

func DBMigrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Note{})
	DB.AutoMigrate(&models.Card{})

}

func UserCheckAvailability(email string) bool {
	var user models.User
	DB.Where("username = ?, email").First(&user)
	return user.ID == 0
}

func CreateUser(user *models.User) (*models.User, error) {

	hashPassword, err := helpers.Encrypt(user)
	if err != nil {
		return nil, err
	}

	isAvailable := CheckUserAvailability(user.Username)
	if isAvailable {
		entry := models.User{Username: user.Username, Password: string(hashPassword)}
		result := DB.Create(&entry)
		return &entry, result.Error
	} else {
		log.Printf("email is not available")
		return nil, nil
	}
}

func UserMatchPassword(username string, password string) (*models.User, error) {
	var user models.User
	DB.Where("username = ?", username).Find(&user)
	if user.ID == 0 {
		return &user, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &models.User{}, nil
	} else {
		return &user, nil
	}
}

func UserFind(id uint64) *models.User {
	var user models.User
	DB.Where("id = ?", id).Find(&user)
	return &user
}

func CreateNote(note *models.Note) (*models.Note, error) {
	entry := models.Note{Name: note.Name, Content: note.Content, UserID: note.UserID}
	result := DB.Create(&entry)
	return &entry, result.Error
}

func CreateCard(card *models.Card) (*models.Card, error) {
	entry := models.Card{Number: card.Number, UserID: card.UserID}
	result := DB.Create(&entry)
	return &entry, result.Error
}

func CheckUserAvailability(username string) bool {
	var user models.User
	DB.Where("username = ?", username).First(&user)
	return user.ID == 0 // true if the username is available
}
