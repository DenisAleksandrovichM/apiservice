package local

import (
	"context"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/cache"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	pb "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"

	"github.com/pkg/errors"
)

const (
	poolSize       = 10
	QuerySortField = "SortField"
	QueryLimit     = "Limit"
	QueryOffset    = "Offset"
	dataBasePort   = ":9081"
)

var (
	errAdd      = errors.New("add error")
	errRead     = errors.New("read error")
	errUpdate   = errors.New("update error")
	errDelete   = errors.New("delete error")
	errList     = errors.New("list error")
	errDataBase = errors.New("database error")
)

func New() cachePkg.Interface {
	return &cache{
		mu:     sync.RWMutex{},
		poolCh: make(chan struct{}, poolSize),
	}
}

type cache struct {
	mu     sync.RWMutex
	poolCh chan struct{}
}

func (c *cache) List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		<-c.poolCh
	}()

	client, err := getDatabaseClient()
	if err != nil {
		return nil, errors.Wrap(errList, err.Error())
	}
	request := &pb.UserListRequest{}
	if val, ok := queryParams[QuerySortField]; ok {
		sortField := val.(string)
		request.SortField = &sortField
	}
	if val, ok := queryParams[QueryLimit]; ok {
		limit := val.(uint64)
		request.Limit = &limit
	}
	if val, ok := queryParams[QueryOffset]; ok {
		offset := val.(uint64)
		request.Offset = &offset
	}
	response, err := client.UserList(ctx, request)
	if err != nil {
		return nil, errors.Wrap(errList, err.Error())
	}
	usersLen := len(response.Users)
	users := make([]models.User, usersLen, usersLen)
	for i, user := range response.Users {
		users[i] = models.User{
			Login:     user.GetLogin(),
			FirstName: user.GetFirstName(),
			LastName:  user.GetLastName(),
			Weight:    float32(user.GetWeight()),
			Height:    uint(user.GetHeight()),
			Age:       uint(user.GetAge()),
		}
	}

	return users, nil
}

func (c *cache) Add(ctx context.Context, user models.User) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	client, err := getDatabaseClient()
	if err != nil {
		return errors.Wrap(errAdd, err.Error())
	}
	_, err = client.UserCreate(ctx, &pb.UserCreateRequest{
		Login:     user.Login,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Weight:    float64(user.Weight),
		Height:    uint32(user.Height),
		Age:       uint32(user.Age),
	})
	if err != nil {
		return errors.Wrap(errAdd, err.Error())
	}

	return nil
}

func (c *cache) Read(ctx context.Context, login string) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	client, err := getDatabaseClient()
	if err != nil {
		return models.User{}, errors.Wrap(errRead, err.Error())
	}
	response, err := client.UserRead(ctx, &pb.UserReadRequest{Login: login})
	if err != nil {
		return models.User{}, errors.Wrap(errRead, err.Error())
	}
	user := models.User{
		Login:     response.Login,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Weight:    float32(response.Weight),
		Height:    uint(response.Height),
		Age:       uint(response.Age),
	}

	return user, nil
}

func (c *cache) Update(ctx context.Context, user models.User) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	client, err := getDatabaseClient()
	if err != nil {
		return errors.Wrap(errUpdate, err.Error())
	}
	_, err = client.UserUpdate(ctx,
		&pb.UserUpdateRequest{
			Login:     user.Login,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Weight:    float64(user.Weight),
			Height:    uint32(user.Height),
			Age:       uint32(user.Age),
		})
	if err != nil {
		return errors.Wrap(errUpdate, err.Error())
	}

	return nil
}

func (c *cache) Delete(ctx context.Context, login string) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	client, err := getDatabaseClient()
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}
	_, err = client.UserDelete(ctx, &pb.UserDeleteRequest{Login: login})
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}

	return nil
}

func getDatabaseClient() (pb.AdminClient, error) {
	conn, err := grpc.Dial(dataBasePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(errDataBase, err.Error())
	}

	return pb.NewAdminClient(conn), nil
}
