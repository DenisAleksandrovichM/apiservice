package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/commander"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/storage"
)

const (
	helpCmd   = "help"
	startCmd  = "start"
	listCmd   = "list"
	addCmd    = "add"
	readCmd   = "read"
	updateCmd = "update"
	deleteCmd = "delete"
)

var BadArgument = errors.New("bad argument")

func listFunc(_ string) (string, error) {
	users := storage.List()
	res := make([]string, 0, len(users))
	for _, user := range users {
		res = append(res, user.String())
	}
	return strings.Join(res, "\n"), nil
}

func startFunc(_ string) (string, error) {
	return "Hello! This bot is designed to display user data.\n" +
		"Each user has fields first name(string), last name(string), weight(float32), height(uint), age(unint).\n" +
		"To work with the bot, you can use the commands:\n" +
		"/help - list commands\n" +
		"/list - list data\n" +
		"/add <first name> <last name> <weight> <height> <age> - add new user with first name, last name, weight, height and age\n" +
		"/delete <id> - delete user by id\n" +
		"/update <id> <first name> <last name> <weight> <height> <age> - update user data by id\n" +
		"/read <id> - reads and displays user by id", nil
}

func helpFunc(_ string) (string, error) {
	return "/help - list commands\n" +
		"/list - list data\n" +
		"/add <first name> <last name> <weight> <height> <age> - add new user with first name, last name, weight, height and age\n" +
		"/delete <id> - delete user by id\n" +
		"/update <id> <first name> <last name> <weight> <height> <age> - update user data by id\n" +
		"/read <id> - reads and displays user by id", nil
}

func addFunc(data string) (string, error) {
	log.Printf("add command param: <data>")
	if data == "" {
		return "", errors.Wrap(BadArgument, "data is empty")
		//return errors.Wrap(BadArgument, "data is empty").Error()
	}
	params := strings.Split(data, " ")
	if len(params) != 5 {
		return "", errors.Wrapf(BadArgument, "expected 5 arguments. You entered %d arguments: %v", len(params), params)
	}
	firstName, lastName, weight, height, age := params[0], params[1], params[2], params[3], params[4]
	user, err := newUser(firstName, lastName, weight, height, age)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	if err = storage.Add(user); err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	return fmt.Sprintf("%v added", user), nil
}

func deleteFunc(id string) (string, error) {
	log.Printf("delete command param: <id>")
	if id == "" {
		return "", errors.Wrap(BadArgument, "id is empty")
	}
	if idUint, err := strconv.ParseUint(id, 0, 0); err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	} else if err = storage.Delete(uint(idUint)); err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	return fmt.Sprintf("user %s deleted", id), nil
}

func updateFunc(data string) (string, error) {
	log.Printf("update command param: <data>")
	if data == "" {
		return "", errors.Wrap(BadArgument, "data is empty")
	}
	params := strings.Split(data, " ")
	if len(params) != 6 {
		return "", errors.Wrapf(BadArgument, "expected 6 arguments. You entered %d arguments: <%v>", len(params), params)
	}
	id, firstName, lastName, weight, height, age := params[0], params[1], params[2], params[3], params[4], params[5]
	user, err := newUser(firstName, lastName, weight, height, age)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	userId, err := strconv.ParseUint(id, 0, 0)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	user.SetId(uint(userId))
	if err = storage.Update(user); err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	return fmt.Sprintf("%v updated", user), nil
}

func readFunc(idStr string) (string, error) {
	if idStr == "" {
		return "", errors.Wrap(BadArgument, "id is empty")
	}
	id, err := strconv.ParseUint(idStr, 0, 0)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	user, err := storage.Read(uint(id))
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	return user.String(), nil
}

func newUser(firstName, lastName, weightStr, heightStr, ageStr string) (*storage.User, error) {
	weight, err := strconv.ParseFloat(weightStr, 32)
	if err != nil {
		return nil, errors.Wrap(BadArgument, err.Error())
	}
	height, err := strconv.ParseUint(heightStr, 0, 0)
	if err != nil {
		return nil, errors.Wrap(BadArgument, err.Error())
	}
	age, err := strconv.ParseUint(ageStr, 0, 0)
	if err != nil {
		return nil, errors.Wrap(BadArgument, err.Error())
	}

	user, err := storage.NewUser(firstName, lastName, float32(weight), uint(height), uint(age))
	if err != nil {
		return nil, errors.Wrap(BadArgument, err.Error())
	}
	return user, nil
}

func AddHandlers(c *commander.Commander) {
	c.RegisterHandler(startCmd, startFunc)
	c.RegisterHandler(helpCmd, helpFunc)
	c.RegisterHandler(listCmd, listFunc)
	c.RegisterHandler(addCmd, addFunc)
	c.RegisterHandler(deleteCmd, deleteFunc)
	c.RegisterHandler(updateCmd, updateFunc)
	c.RegisterHandler(readCmd, readFunc)
}
