package local

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/cache"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	pb "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

const (
	poolSize       = 10
	QuerySortField = "SortField"
	QueryLimit     = "Limit"
	QueryOffset    = "Offset"
	dataBasePort   = ":9011"
	topic          = "test_2808"
	eventCreate    = "Create"
	eventUpdate    = "Update"
)

var (
	errAdd      = errors.New("add error")
	errRead     = errors.New("read error")
	errUpdate   = errors.New("update error")
	errDelete   = errors.New("delete error")
	errList     = errors.New("list error")
	errDataBase = errors.New("database error")
	errKafka    = errors.New("kafka error")
	brokers     = []string{"localhost:19091"}
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

func (c *cache) Add(ctx context.Context, user models.User) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	userJSON, err := json.Marshal(user)
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}
	if err = sendMessage(eventCreate, userJSON); err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}

	time.Sleep(time.Second * 10)

	responseUser, err := getUserByLogin(ctx, user.Login)
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}

	return responseUser, nil
}

func (c *cache) Read(ctx context.Context, login string) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()
	return getUserByLogin(ctx, login)
}

func (c *cache) Update(ctx context.Context, user models.User) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	userJSON, err := json.Marshal(user)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}
	if err = sendMessage(eventUpdate, userJSON); err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}

	time.Sleep(time.Second * 10)

	responseUser, err := getUserByLogin(ctx, user.Login)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}

	return responseUser, errors.Wrap(errUpdate, err.Error())
}

func (c *cache) Delete(_ context.Context, login string) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()
	if err := sendMessage("Delete", []byte(login)); err != nil {
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

func sendMessage(key string, value []byte) error {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return errors.Wrap(errKafka, err.Error())
	}
	_, _, err = syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		return errors.Wrap(errKafka, err.Error())
	}
	return nil
}

func getUserByLogin(ctx context.Context, login string) (models.User, error) {
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
