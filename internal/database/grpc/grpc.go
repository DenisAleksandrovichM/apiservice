package grpc

import (
	"context"
	userPkg "github.com/DenisAleksandrovichM/apiservice/internal/database/core/user"
	pb "github.com/DenisAleksandrovichM/apiservice/pkg/api"
	"github.com/DenisAleksandrovichM/apiservice/pkg/counter/errorsCounter"
	"github.com/DenisAleksandrovichM/apiservice/pkg/counter/requestsCounter"
	"github.com/DenisAleksandrovichM/apiservice/pkg/counter/responsesCounter"
	"github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(user userPkg.User) *implementation {
	return &implementation{
		user: user,
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	user userPkg.User
}

func (i *implementation) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserCreateResponse{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age)}, nil
}

func (i *implementation) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	requestsCounter.Inc()
	var queryParams = map[string]interface{}{}
	if in.Offset != nil {
		queryParams[userPkg.QueryOffset] = in.GetOffset()
	}
	if in.Limit != nil {
		queryParams[userPkg.QueryLimit] = in.GetLimit()
	}
	if in.SortField != nil {
		queryParams[userPkg.QuerySortField] = in.GetSortField()
	}

	users, err := i.user.List(ctx, queryParams)
	if err != nil {
		errorsCounter.Inc()
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

func (i *implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
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
		errorsCounter.Inc()
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserUpdateResponse{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age)}, nil
}

func (i *implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	requestsCounter.Inc()
	err := i.user.Delete(ctx, in.GetLogin())
	if err != nil {
		errorsCounter.Inc()
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsesCounter.Inc()
	return &pb.UserDeleteResponse{}, nil
}

func (i *implementation) UserRead(ctx context.Context, in *pb.UserReadRequest) (*pb.UserReadResponse, error) {
	requestsCounter.Inc()
	user, err := i.user.Read(ctx, in.GetLogin())
	if err != nil {
		errorsCounter.Inc()
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
