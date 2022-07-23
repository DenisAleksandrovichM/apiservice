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

func listFunc(_ string) string {
	users := storage.List()
	res := make([]string, 0, len(users))
	for _, user := range users {
		res = append(res, user.String())
	}
	return strings.Join(res, "\n")
}

func startFunc(_ string) string {
	return "Hello! This bot is designed to display user data.\n" +
		"Each user has fields first name(string), last name(string), weight(float32), height(uint), age(unint).\n" +
		"To work with the bot, you can use the commands:\n" +
		helpFunc("")
}

func helpFunc(_ string) string {
	return "/help - list commands\n" +
		"/list - list data\n" +
		"/add <first name> <last name> <weight> <height> <age> - add new user with first name, last name, weight, height and age\n" +
		"/delete <id> - delete user by id\n" +
		"/update <id> <first name> <last name> <weight> <height> <age> - update user data by id\n" +
		"/read <id> - reads and displays user by id"
}

func addFunc(data string) string {
	log.Printf("add command param: <data>")
	if data == "" {
		return errors.Wrap(BadArgument, "data is empty").Error()
	}
	params := strings.Split(data, " ")
	if len(params) != 5 {
		return errors.Wrapf(BadArgument, "expected 5 arguments. You entered %d arguments: %v", len(params), params).Error()
	}
	user, err := newUser(params[0], params[1], params[2], params[3], params[4])
	if err != nil {
		return errors.Wrap(BadArgument, err.Error()).Error()
	}
	if err = storage.Add(user); err != nil {
		return errors.Wrap(BadArgument, err.Error()).Error()
	}
	return fmt.Sprintf("%v added", user)
}

func deleteFunc(id string) string {
	log.Printf("delete command param: <id>")
	if id == "" {
		return errors.Wrap(BadArgument, "id is empty").Error()
	}
	if idUint, err := strconv.ParseUint(id, 0, 0); err != nil {
		return err.Error()
	} else if err = storage.Delete(uint(idUint)); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("user %s deleted", id)
}

func updateFunc(data string) string {
	log.Printf("update command param: <data>")
	if data == "" {
		return errors.Wrap(BadArgument, "data is empty").Error()
	}
	params := strings.Split(data, " ")
	if len(params) != 6 {
		return errors.Wrapf(BadArgument, "expected 6 arguments. You entered %d arguments: <%v>", len(params), params).Error()
	}
	user, err := newUser(params[1], params[2], params[3], params[4], params[5])
	if err != nil {
		return errors.Wrap(BadArgument, err.Error()).Error()
	}
	if id, err := strconv.ParseUint(params[0], 0, 0); err != nil {
		return errors.Wrap(BadArgument, err.Error()).Error()
	} else {
		user.SetId(uint(id))
	}
	if err = storage.Update(user); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v updated", user)
}

func readFunc(idStr string) string {
	if idStr == "" {
		return errors.Wrap(BadArgument, "id is empty").Error()
	}
	if id, err := strconv.ParseUint(idStr, 0, 0); err == nil {
		if user, err := storage.Read(uint(id)); err == nil {
			return user.String()
		} else {
			return errors.Wrap(BadArgument, err.Error()).Error()
		}
	} else {
		return errors.Wrap(BadArgument, err.Error()).Error()
	}
}

func newUser(firstName, lastName, weightStr, heightStr, ageStr string) (*storage.User, error) {
	weight, err := strconv.ParseFloat(weightStr, 32)
	if err != nil {
		return nil, err
	}
	height, err := strconv.ParseUint(heightStr, 0, 0)
	if err != nil {
		return nil, err
	}
	age, err := strconv.ParseUint(ageStr, 0, 0)
	if err != nil {
		return nil, err
	}

	if user, err := storage.NewUser(firstName, lastName, float32(weight), uint(height), uint(age)); err == nil {
		return user, nil
	} else {
		return nil, err
	}
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
