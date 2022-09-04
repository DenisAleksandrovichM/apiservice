package local

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/cache"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/hitCounter"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/counter/missCounter"
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
	dataBasePort   = ":9091"
	topic          = "test_2808"
	eventCreate    = "Create"
	eventUpdate    = "Update"
	eventDelete    = "Delete"
	redisAddress   = "localhost:6379"
	redisDB        = 0
	redisPassword  = ""
	waitingtime    = time.Second * 10
	emptyID        = ""
)

var (
	brokers              = []string{"localhost:19091"}
	errAdd               = errors.New("add error")
	errRead              = errors.New("read error")
	errUpdate            = errors.New("update error")
	errDelete            = errors.New("delete error")
	errList              = errors.New("list error")
	errDataBase          = errors.New("database error")
	errKafka             = errors.New("kafka error")
	errRedisGet          = errors.New("redis on get error")
	errRedisScan         = errors.New("redis on scan error")
	errCheckCUResult     = errors.New("check create/update result error")
	errCheckDeleteResult = errors.New("check delete result error")
	errMarshal           = errors.New("on marshal error")
)

func New() cachePkg.Interface {
	return &cache{
		mu:     sync.RWMutex{},
		poolCh: make(chan struct{}, poolSize),
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			DB:       redisDB,
			Password: redisPassword,
		}),
	}
}

type cache struct {
	mu          sync.RWMutex
	poolCh      chan struct{}
	redisClient *redis.Client
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

func (c *cache) Add(_ context.Context, user models.User) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	id, err := sendMessage(eventCreate, user)
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}
	responseUser, err := c.checkCUDResult(id)
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
	return c.getUserByLogin(ctx, login)
}

func (c *cache) Update(_ context.Context, user models.User) (models.User, error) {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()
	id, err := sendMessage(eventUpdate, user)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}
	responseUser, err := c.checkCUDResult(id)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}
	return responseUser, nil
}

func (c *cache) Delete(_ context.Context, login string) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()
	id, err := sendMessage(eventDelete, login)
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}
	_, err = c.checkCUDResult(id)
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

func (c *cache) getUserByLogin(ctx context.Context, login string) (models.User, error) {
	user, err := c.getRedisUser(login)
	if err != nil {
		log.Info(err)
		client, err := getDatabaseClient()
		if err != nil {
			return models.User{}, errors.Wrap(errRead, err.Error())
		}
		response, err := client.UserRead(ctx, &pb.UserReadRequest{Login: login})
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

func (c *cache) getRedisUser(key string) (models.User, error) {
	getUser := c.redisClient.Get(key)
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

func (c *cache) checkCUDResult(id string) (models.User, error) {
	time.Sleep(waitingtime)
	response, err := c.getRedisUser(id)
	if err != nil {
		return models.User{}, errors.Wrap(errCheckCUResult, err.Error())
	}
	return response, nil
}
