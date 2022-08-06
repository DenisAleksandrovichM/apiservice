package api

import (
	"context"
	"github.com/pkg/errors"
	userPkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/cache/local"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/validate"
	pb "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(user userPkg.Interface) pb.AdminServer {
	return &implementation{
		user: user,
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	user userPkg.Interface
}

func (i *implementation) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	if err := i.user.Create(ctx, models.User{
		Login:     in.GetLogin(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Weight:    float32(in.GetWeight()),
		Height:    uint(in.GetHeight()),
		Age:       uint(in.GetAge()),
	}); err != nil {
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserCreateResponse{}, nil
}

func (i *implementation) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {

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

	return &pb.UserListResponse{
		Users: result,
	}, nil
}

func (i *implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	if err := i.user.Update(ctx, models.User{
		Login:     in.GetLogin(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Weight:    float32(in.GetWeight()),
		Height:    uint(in.GetHeight()),
		Age:       uint(in.GetAge()),
	}); err != nil {
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserUpdateResponse{}, nil
}

func (i *implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	if err := i.user.Delete(ctx, in.GetLogin()); err != nil {
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserDeleteResponse{}, nil
}

func (i *implementation) UserRead(ctx context.Context, in *pb.UserReadRequest) (*pb.UserReadResponse, error) {
	user, err := i.user.Read(ctx, in.GetLogin())
	if err != nil {
		if errors.Is(err, validate.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserReadResponse{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age)}, nil
}
