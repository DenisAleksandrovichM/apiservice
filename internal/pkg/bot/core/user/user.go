//go:generate mockgen -source ./user.go -destination=./mocks/user.go -package=mock_user
package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	"github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/validate"
	"github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/hitCounter"
	"github.com/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/missCounter"
	pb "github.com/DenisAleksandrovichM/homework-1/pkg/api"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

var brokers = []string{"localhost:19091"}

type User interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Read(ctx context.Context, login string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, login string) error
	List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error)
	String(user models.User) string
}

type implementation struct {
	mu          sync.RWMutex
	poolCh      chan struct{}
	redisClient *redis.Client
	dbClient    pb.AdminClient
}

func New() *implementation {
	dbClient, err := getDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}
	return &implementation{
		mu:     sync.RWMutex{},
		poolCh: make(chan struct{}, poolSize),
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			DB:       redisDB,
			Password: redisPassword,
		},
		),
		dbClient: dbClient,
	}
}

func (i *implementation) Create(_ context.Context, user models.User) (models.User, error) {
	if err := validate.ValidateUser(user); err != nil {
		return models.User{}, err
	}
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()

	id, err := sendMessage(eventCreate, user)
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}
	responseUser, err := i.checkCUDResult(id)
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}
	return responseUser, nil
}

func (i *implementation) Read(ctx context.Context, login string) (models.User, error) {
	if err := validate.ValidateLogin(login); err != nil {
		return models.User{}, err
	}
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()
	return i.getUserByLogin(ctx, login)
}

func (i *implementation) Update(_ context.Context, user models.User) (models.User, error) {
	if err := validate.ValidateUser(user); err != nil {
		return models.User{}, err
	}
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()
	id, err := sendMessage(eventUpdate, user)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}
	responseUser, err := i.checkCUDResult(id)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}
	return responseUser, nil
}

func (i *implementation) Delete(_ context.Context, login string) error {
	if err := validate.ValidateLogin(login); err != nil {
		return err
	}
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()
	id, err := sendMessage(eventDelete, login)
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}
	_, err = i.checkCUDResult(id)
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}
	return nil
}

func (i *implementation) List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error) {
	i.poolCh <- struct{}{}
	i.mu.RLock()
	defer func() {
		i.mu.RUnlock()
		<-i.poolCh
	}()

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
	response, err := i.dbClient.UserList(ctx, request)
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

func (i *implementation) String(user models.User) string {
	return fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
		user.Login, user.FirstName, user.LastName, user.Weight, user.Height, user.Age)
}

func (i *implementation) checkCUDResult(id string) (models.User, error) {
	time.Sleep(waitingtime)
	response, err := i.getRedisUser(id)
	if err != nil {
		return models.User{}, errors.Wrap(errCheckCUResult, err.Error())
	}
	return response, nil
}

func (i *implementation) getRedisUser(key string) (models.User, error) {
	getUser := i.redisClient.Get(key)
	if getUser.Err() != nil {
		missCounter.Inc()
		return models.User{}, errors.Wrap(errRedisGet, getUser.Err().Error())
	}
	hitCounter.Inc()
	user := models.User{}
	err := getUser.Scan(&user)
	if err != nil {
		return models.User{}, errors.Wrap(errRedisScan, err.Error())
	}
	return user, nil
}

func (i *implementation) getUserByLogin(ctx context.Context, login string) (models.User, error) {
	user, err := i.getRedisUser(login)
	if err != nil {
		log.Info(err)
		response, err := i.dbClient.UserRead(ctx, &pb.UserReadRequest{Login: login})
		if err != nil {
			return models.User{}, errors.Wrap(errRead, err.Error())
		}
		user = models.User{
			Login:     response.Login,
			FirstName: response.FirstName,
			LastName:  response.LastName,
			Weight:    float32(response.Weight),
			Height:    uint(response.Height),
			Age:       uint(response.Age),
		}
	}
	return user, nil
}

func getDatabaseClient() (pb.AdminClient, error) {
	conn, err := grpc.Dial(dataBasePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(errDataBase, err.Error())
	}

	return pb.NewAdminClient(conn), nil
}

func sendMessage(key string, value interface{}) (string, error) {
	msg, err := newKafkaMessage(value)
	if err != nil {
		return emptyID, errors.Wrap(errKafka, err.Error())
	}
	msgBin, err := msg.MarshalBinary()
	if err != nil {
		return emptyID, errors.Wrap(errKafka, err.Error())
	}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return emptyID, errors.Wrap(errKafka, err.Error())
	}
	_, _, err = syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msgBin),
	})
	if err != nil {
		return emptyID, errors.Wrap(errKafka, err.Error())
	}
	return msg.Id, nil
}

func newKafkaMessage(data interface{}) (*models.KafkaMessage, error) {
	var dataMsg []byte
	if dataStr, ok := data.(string); ok {
		dataMsg = []byte(dataStr)
	} else {
		msg, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(errMarshal, err.Error())
		}
		dataMsg = msg
	}
	return &models.KafkaMessage{
		Id:   uuid.New().String(),
		Data: dataMsg,
	}, nil
}
