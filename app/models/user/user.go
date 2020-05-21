package user

import (
	"bacancy/go-boiler-plate/app/common"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Password     string `gorm:"not null" json:"-"`
	Salt         string `gorm:"not null" json:"-"`
	Email        string `gorm:"not null;unique"`
	Admin        bool   `gorm:"not null"`
	RecoveryCode string `gorm:"null" json:"-"`
}

func Create(email string, password string, salt string, name string, lastname string) (*User, error) {

	user := User{
		Name:     name,
		LastName: lastname,
		Password: password,
		Salt:     salt,
		Email:    email,
		Admin:    false,
	}

	err := common.GetDatabase().Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserById(id uint) (user User, found bool, err error) {
	user = User{}

	r := common.GetDatabase()

	r = r.Where("id = ?", int(id)).First(&user)
	if r.RecordNotFound() {
		return user, false, r.Error
	}

	if r.Error != nil {
		return user, true, r.Error
	}

	return user, true, nil
}

func GetUserByEmail(email string) (user User, found bool, err error) {
	user = User{}

	r := common.GetDatabase()

	r = r.Where("email = ?", email).First(&user)
	if r.RecordNotFound() {
		return user, false, nil
	}

	if r.Error != nil {
		return user, true, r.Error
	}

	return user, true, nil
}

func ChangePassword(id uint, newPassword string, salt string) error {
	db := common.GetDatabase()

	update := db.Table("users").Where("id = ?", id).Updates(map[string]interface{}{"Password": newPassword, "recovery_code": 0, "salt": salt})
	if update.Error != nil {
		return update.Error
	}

	return nil
}

func ChangeProfileData(id uint, name string, phone string, lastname string) error {

	db := common.GetDatabase()

	update := db.Table("users").Where("id = ?", id).Updates(map[string]interface{}{"name": name, "phone": phone, "last_name": lastname})
	if update.Error != nil {
		return update.Error
	}

	return nil
}

func GetUserProfile(id uint) (user User, found bool, err error) {

	user = User{}

	r := common.GetDatabase()

	r = r.Where("id = ?", id).First(&user)
	if r.RecordNotFound() {
		return user, false, r.Error
	}

	if r.Error != nil {
		return user, true, r.Error
	}

	return user, true, nil

}
