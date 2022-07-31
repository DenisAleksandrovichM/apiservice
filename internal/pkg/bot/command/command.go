package command

import (
	"github.com/pkg/errors"
	modelsPkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/models"
	validatePkg "gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/validate"
	"strconv"
	"strings"
)

var BadArgument = errors.New("bad argument")

type Interface interface {
	Name() string
	Description() string
	Process(args string) (string, error)
}

func ProcessAddOrUpdate(args string, function func(user modelsPkg.User) error) (string, error) {
	params, err := ValidateParams(args, 6)
	if err != nil {
		return "", err
	}
	login, firstName, lastName, weightStr, heightStr, ageStr := params[0], params[1], params[2], params[3], params[4], params[5]
	weight, err := strconv.ParseFloat(weightStr, 32)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	height, err := strconv.ParseUint(heightStr, 0, 0)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}
	age, err := strconv.ParseUint(ageStr, 0, 0)
	if err != nil {
		return "", errors.Wrap(BadArgument, err.Error())
	}

	if err = function(modelsPkg.User{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Weight:    float32(weight),
		Height:    uint(height),
		Age:       uint(age),
	}); err != nil {
		if errors.Is(err, validatePkg.ErrValidation) {
			return "", errors.Wrap(BadArgument, err.Error())
		}
		return "", errors.Wrap(BadArgument, "internal error")
	}

	return "success", nil
}

func ValidateParams(args string, numOfParams int) ([]string, error) {
	params := strings.Split(args, " ")
	if len(params) != numOfParams {
		return nil, errors.Wrapf(BadArgument, "expected %d arguments. You entered %d arguments: %v", numOfParams, len(params), params)
	}
	return params, nil
}
