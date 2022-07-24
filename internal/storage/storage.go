package storage

import (
	"log"
	"strconv"

	"github.com/pkg/errors"
)

var data map[uint]*User

var UserNotExists = errors.New("user does not exist")
var UserExists = errors.New("user exists")

func init() {
	log.Println("init storage")
	data = make(map[uint]*User)
	u, err := NewUser("Denis", "Merzlikin", 90.5, 190, 29)
	if err != nil {
		log.Fatal(err)
	}
	if err = Add(u); err != nil {
		log.Fatal(err)
	}
}

func List() []*User {
	res := make([]*User, 0, len(data))
	for _, v := range data {
		res = append(res, v)
	}
	return res
}

func Add(user *User) error {
	if _, ok := data[user.GetId()]; ok {
		return errors.Wrap(UserExists, strconv.FormatUint(uint64(user.GetId()), 10))
	}
	data[user.GetId()] = user
	return nil
}

func Update(user *User) error {
	if _, ok := data[user.GetId()]; !ok {
		return errors.Wrap(UserNotExists, strconv.FormatUint(uint64(user.GetId()), 10))
	}
	data[user.GetId()] = user
	return nil
}

func Delete(id uint) error {
	if _, ok := data[id]; ok {
		delete(data, id)
		return nil
	}
	return errors.Wrap(UserNotExists, strconv.FormatUint(uint64(id), 10))
}

func Read(id uint) (*User, error) {
	if user, ok := data[id]; ok {
		return user, nil
	}
	return nil, errors.Wrap(UserNotExists, strconv.FormatUint(uint64(id), 10))
}
