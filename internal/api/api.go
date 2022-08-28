package api

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/cache/local"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/validate"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/errorsCounter"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/requestsCounter"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/responsesCounter"
	pb "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Implementation struct {
	pb.UnimplementedAdminServer
	user userPkg.Interface
}

func New(user userPkg.Interface) *Implementation {
	return &Implementation{
		user: user,
	}
}

func (i *Implementation) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	requestsCounter.Inc()
	user, err := i.user.Create(ctx, models.User{
		Login:     in.GetLogin(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Weight:    float32(in.GetWeight()),
		Height:    uint(in.GetHeight()),
		Age:       uint(in.GetAge()),
	})
	if err != nil {
		errorsCounter.Inc()
		log.Error(err)
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserCreateResponse{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age),
	}, nil
}

func (i *Implementation) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	requestsCounter.Inc()
	var queryParams = map[string]interface{}{}
	if in.Offset != nil {
		queryParams[local.QueryOffset] = in.GetOffset()
	}
	if in.Limit != nil {
		queryParams[local.QueryLimit] = in.GetLimit()
	}
	if in.SortField != nil {
		queryParams[local.QuerySortField] = in.GetSortField()
	}

	users, err := i.user.List(ctx, queryParams)
	if err != nil {
		errorsCounter.Inc()
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*pb.UserListResponse_User, 0, len(users))
	for _, user := range users {
		result = append(result, &pb.UserListResponse_User{
			Login:     user.Login,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Weight:    float64(user.Weight),
			Height:    uint32(user.Height),
			Age:       uint32(user.Age),
		})
	}

	responsesCounter.Inc()
	return &pb.UserListResponse{
		Users: result,
	}, nil
}

func (i *Implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	requestsCounter.Inc()
	user, err := i.user.Update(ctx, models.User{
		Login:     in.GetLogin(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Weight:    float32(in.GetWeight()),
		Height:    uint(in.GetHeight()),
		Age:       uint(in.GetAge()),
	})
	if err != nil {
		log.Error(err)
		errorsCounter.Inc()
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserUpdateResponse{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age),
	}, nil
}

func (i *Implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	requestsCounter.Inc()
	err := i.user.Delete(ctx, in.GetLogin())
	if err != nil {
		errorsCounter.Inc()
		log.Error(err)
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserDeleteResponse{}, nil
}

func (i *Implementation) UserRead(ctx context.Context, in *pb.UserReadRequest) (*pb.UserReadResponse, error) {
	requestsCounter.Inc()
	user, err := i.user.Read(ctx, in.GetLogin())
	if err != nil {
		errorsCounter.Inc()
		log.Error(err)
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserReadResponse{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age)}, nil
}
