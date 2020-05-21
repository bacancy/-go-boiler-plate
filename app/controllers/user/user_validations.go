package user

import (
	"bacancy/go-boiler-plate/app/models/user"
	"bacancy/go-boiler-plate/app/validators"
	"errors"
)

func Validate(password string, email string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	if err := validateEmail(email); err != nil {
		return err
	}

	return nil
}

func validatePassword(password string) error {
	if !validators.IsLongerOrEqualThan(password, 6) {
		return errors.New("Password must have 6 or more characters")
	}
	if !validators.IsShorterOrEqualThan(password, 63) {
		return errors.New("Password must have less than 63 characters")
	}

	return nil
}

func validateEmail(email string) error {
	if validators.IsEmpty(email) {
		return errors.New("Email address not informed")
	}
	if !validators.IsEmail(email) {
		return errors.New("Invalid email address")
	}
	_, found, _ := user.GetUserByEmail(email)
	if found {
		return errors.New("Email already in use")
	}

	return nil
}

func validateRecoveryData(email string, recoveryCode string) (uint, error) {

	userData, found, err := user.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}
	if found == false {
		return 0, errors.New("Email not registered")
	}

	if userData, _, _ := user.GetUserById(userData.ID); userData.RecoveryCode != recoveryCode {
		return 0, errors.New("Wrong recovery code")
	}

	return userData.ID, nil
}

func validateLiteralEmail(email string) error {
	if validators.IsEmpty(email) {
		return errors.New("Email address not informed")
	}
	if !validators.IsEmail(email) {
		return errors.New("Invalid email")
	}

	return nil
}
