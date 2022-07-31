package validate

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/models"
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
	if !isValidPropertyValue(float32(len(login)), 0, 15) {
		return errors.Wrapf(ErrValidation, "field: [login] <%s>(maximum 15 characters)", login)
	}
	return nil
}

func validateFirstName(firstName string) error {
	if !isValidPropertyValue(float32(len(firstName)), 0, 10) {
		return errors.Wrapf(ErrValidation, "field: [first name] <%s>(maximum 10 characters)", firstName)
	}
	return nil
}

func validateLastName(lastName string) error {
	if !isValidPropertyValue(float32(len(lastName)), 0, 20) {
		return errors.Wrapf(ErrValidation, "field: [last name] <%s>(maximum 20 characters)", lastName)
	}
	return nil
}

func validateWeight(weight float32) error {
	if !isValidPropertyValue(weight, 20, 250) {
		return errors.Wrapf(ErrValidation, "field: [weight] <%.2f>(min 20, max 250)", weight)
	}
	return nil
}

func validateHeight(height uint) error {
	if !isValidPropertyValue(float32(height), 20, 250) {
		return errors.Wrapf(ErrValidation, "field: [weight] <%d>(min 20, max 250)", height)
	}
	return nil
}

func validateAge(age uint) error {
	if !isValidPropertyValue(float32(age), 0, 150) {
		return errors.Wrapf(ErrValidation, "field: [age] <%d>(max 150)", age)
	}
	return nil
}
