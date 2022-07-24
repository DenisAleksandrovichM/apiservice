package storage

import (
	"fmt"
	"github.com/pkg/errors"
)

var lastId = uint(0)

type User struct {
	id        uint
	firstName string
	lastName  string
	weight    float32
	height    uint
	age       uint
}

func NewUser(firstName, lastName string, weight float32, height, age uint) (*User, error) {
	user := User{}
	errBadArgument := errors.New("bad argument")
	if err := user.SetFirstName(firstName); err != nil {
		return nil, errors.Wrap(errBadArgument, err.Error())
	}
	if err := user.SetLastName(lastName); err != nil {
		return nil, errors.Wrap(errBadArgument, err.Error())
	}
	if err := user.SetWeight(weight); err != nil {
		return nil, errors.Wrap(errBadArgument, err.Error())
	}
	if err := user.SetHeight(height); err != nil {
		return nil, errors.Wrap(errBadArgument, err.Error())
	}
	if err := user.SetAge(age); err != nil {
		return nil, errors.Wrap(errBadArgument, err.Error())
	}
	lastId++
	user.SetId(lastId)
	return &user, nil
}

func (u *User) GetId() uint {
	return u.id
}

func (u *User) SetId(id uint) {
	u.id = id
}

func (u *User) GetFirstName() string {
	return u.firstName
}

func (u *User) SetFirstName(firstName string) error {
	if !isValidPropertyValue(float32(len(firstName)), 0, 10) {
		return fmt.Errorf("bad first name <%s>(maximum 10 characters)", firstName)
	}
	u.firstName = firstName
	return nil
}

func (u *User) GetLastName() string {
	return u.lastName
}

func (u *User) SetLastName(lastName string) error {
	if !isValidPropertyValue(float32(len(lastName)), 0, 20) {
		return fmt.Errorf("bad last name <%s>(maximum 20 characters)", lastName)
	}
	u.lastName = lastName
	return nil
}

func (u *User) GetWeight() float32 {
	return u.weight
}

func (u *User) SetWeight(weight float32) error {
	if !isValidPropertyValue(weight, 20, 250) {
		return fmt.Errorf("bad weight <%.2f>(min 20, max 250)", weight)
	}
	u.weight = weight
	return nil
}

func (u *User) GetHeight() uint {
	return u.height
}

func (u *User) SetHeight(height uint) error {
	if !isValidPropertyValue(float32(height), 20, 250) {
		return fmt.Errorf("bad height <%d>(min 20, max 250)", height)
	}
	u.height = height
	return nil
}

func (u *User) GetAge() uint {
	return u.age
}

func (u *User) SetAge(age uint) error {
	if !isValidPropertyValue(float32(age), 0, 150) {
		return fmt.Errorf("bad age <%d>(max 150)", age)
	}
	u.age = age
	return nil
}

func (u User) String() string {
	return fmt.Sprintf("User %d: %s %s / weight: %0.2f / height: %d / age: %d",
		u.GetId(), u.GetFirstName(), u.GetLastName(), u.GetWeight(), u.GetHeight(), u.GetAge())
}

func isValidPropertyValue(num, min, max float32) bool {
	if num >= min && num <= max {
		return true
	}
	return false
}
