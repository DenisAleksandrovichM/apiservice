package validate

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
)

const (
	minLenLogin     = 1
	maxLenLogin     = 15
	minLenFirstName = 1
	maxLenFirstName = 20
	minLenLastName  = 1
	maxLenLastName  = 25
	minWeight       = 3
	maxWeight       = 250
	minHeight       = 50
	maxHeight       = 250
	minAge          = 0
	maxAge          = 150
)

var ErrValidation = errors.New("invalid data")

func ValidateUser(user models.User) error {
	if err := ValidateLogin(user.Login); err != nil {
		return errors.Wrapf(ErrValidation, err.Error())
	}
	if err := validateFirstName(user.FirstName); err != nil {
		return errors.Wrapf(ErrValidation, err.Error())
	}
	if err := validateLastName(user.LastName); err != nil {
		return errors.Wrapf(ErrValidation, err.Error())
	}
	if err := validateWeight(user.Weight); err != nil {
		return errors.Wrapf(ErrValidation, err.Error())
	}
	if err := validateHeight(user.Height); err != nil {
		return errors.Wrapf(ErrValidation, err.Error())
	}
	if err := validateAge(user.Age); err != nil {
		return errors.Wrapf(ErrValidation, err.Error())
	}
	return nil

}

func isValidPropertyValue(num, min, max float32) bool {
	if num >= min && num <= max {
		return true
	}
	return false
}

func ValidateLogin(login string) error {
	if !isValidPropertyValue(float32(len(login)), minLenLogin, maxLenLogin) {
		return errors.Wrapf(ErrValidation,
			"field: [login] <%s>(min %d, max %d characters)", login, minLenLogin, maxLenLogin)
	}
	return nil
}

func validateFirstName(firstName string) error {
	if !isValidPropertyValue(float32(len(firstName)), minLenFirstName, maxLenFirstName) {
		return errors.Wrapf(ErrValidation,
			"field: [first name] <%s>(min %d, max %d characters)", firstName, minLenFirstName, maxLenFirstName)
	}
	return nil
}

func validateLastName(lastName string) error {
	if !isValidPropertyValue(float32(len(lastName)), minLenLastName, maxLenLastName) {
		return errors.Wrapf(ErrValidation,
			"field: [last name] <%s>(min %d, max %d characters)", lastName, minLenLastName, maxLenLastName)
	}
	return nil
}

func validateWeight(weight float32) error {
	if !isValidPropertyValue(weight, minWeight, maxWeight) {
		return errors.Wrapf(ErrValidation, "field: [weight] <%.2f>(min %d, max %d)", weight, minWeight, maxWeight)
	}
	return nil
}

func validateHeight(height uint) error {
	if !isValidPropertyValue(float32(height), minHeight, maxHeight) {
		return errors.Wrapf(ErrValidation, "field: [height] <%d>(min %d, max %d)", height, minHeight, maxHeight)
	}
	return nil
}

func validateAge(age uint) error {
	if !isValidPropertyValue(float32(age), minAge, maxAge) {
		return errors.Wrapf(ErrValidation, "field: [age] <%d>(max %d)", age, maxAge)
	}
	return nil
}
