//go:generate mockgen -source ./command.go -destination=./mocks/command.go -package=mock_command
package command

import (
	"context"
	validatePkg "github.com/DenisAleksandrovichM/apiservice/internal/api/core/user/validate"
	modelsPkg "github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var BadArgument = errors.New("bad argument")

type Command interface {
	Name() string
	Description() string
	Process(ctx context.Context, args string) (string, error)
}

func ProcessAddOrUpdate(ctx context.Context, args string, function func(context.Context, modelsPkg.User) (modelsPkg.User, error)) (modelsPkg.User, error) {
	params, err := ValidateParams(args, 6)
	if err != nil {
		return modelsPkg.User{}, errors.Wrap(BadArgument, err.Error())
	}
	login, firstName, lastName, weightStr, heightStr, ageStr := params[0], params[1], params[2], params[3], params[4], params[5]
	weight, err := strconv.ParseFloat(weightStr, 32)
	if err != nil {
		return modelsPkg.User{}, errors.Wrap(BadArgument, err.Error())
	}
	height, err := strconv.ParseUint(heightStr, 0, 0)
	if err != nil {
		return modelsPkg.User{}, errors.Wrap(BadArgument, err.Error())
	}
	age, err := strconv.ParseUint(ageStr, 0, 0)
	if err != nil {
		return modelsPkg.User{}, errors.Wrap(BadArgument, err.Error())
	}
	user, err := function(ctx,
		modelsPkg.User{
			Login:     strings.ToLower(login),
			FirstName: firstName,
			LastName:  lastName,
			Weight:    float32(weight),
			Height:    uint(height),
			Age:       uint(age),
		})
	if err != nil {
		if errors.Is(err, validatePkg.ErrValidation) {
			return modelsPkg.User{}, errors.Wrap(BadArgument, err.Error())
		}
		return modelsPkg.User{}, errors.Wrap(BadArgument, "internal error")
	}

	return user, nil
}

func ValidateParams(args string, numOfParams int) ([]string, error) {
	params := strings.Split(args, " ")
	if len(params) != numOfParams {
		return nil, errors.Wrapf(BadArgument,
			"expected %d arguments. You entered %d arguments: %v", numOfParams, len(params), params)
	}
	return params, nil
}
