//go:generate mockgen -source ./user.go -destination=./mocks/user.go -package=mock_user
package user

import (
	"context"
	"fmt"
	"github.com/DenisAleksandrovichM/apiservice/pkg/counter/hitCounter"
	"github.com/DenisAleksandrovichM/apiservice/pkg/counter/missCounter"
	"github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type User interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Read(ctx context.Context, login string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, login string) error
	List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error)
	String(user models.User) string
	SetRedisData(key string, value interface{})
}

func New(db *sqlx.DB) User {
	return &implementation{
		mu:     sync.RWMutex{},
		poolCh: make(chan struct{}, poolSize),
		db:     db,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			DB:       redisDB,
			Password: redisPassword,
		}),
	}
}

type implementation struct {
	mu          sync.RWMutex
	poolCh      chan struct{}
	db          *sqlx.DB
	redisClient *redis.Client
}

func (i *implementation) Create(ctx context.Context, user models.User) (models.User, error) {
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()

	_, err := i.getUserByLogin(user.Login)
	if err != nil && !errors.Is(err, errUserNotExists) {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	} else if err == nil {
		return models.User{}, errors.Wrap(errAdd, errUserExists.Error())
	}

	query, args, err := squirrel.Insert(tableName).
		Columns(tableColumns).
		Values(strings.ToLower(user.Login), user.FirstName, user.LastName, user.Weight, user.Height, user.Age).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumns)).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}

	queryer := i.getQueryer(nil)
	row := queryer.QueryRowxContext(ctx, query, args...)

	var (
		login, firstName, lastName string
		weight                     float32
		height, age                uint
	)

	err = row.Scan(&login, &firstName, &lastName, &weight, &height, &age)
	if err != nil {
		return models.User{}, errors.Wrap(errAdd, err.Error())
	}
	createdUser := models.User{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Weight:    weight,
		Height:    height,
		Age:       age}
	return createdUser, nil
}

func (i *implementation) Read(_ context.Context, login string) (models.User, error) {
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()

	user, err := i.getUserByLogin(login)
	if err != nil {
		return models.User{}, errors.Wrap(errRead, err.Error())
	}

	return user, nil
}

func (i *implementation) Update(ctx context.Context, user models.User) (models.User, error) {
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()

	_, err := i.getUserByLogin(user.Login)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}

	setMap := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"weight":     user.Weight,
		"height":     user.Height,
		"age":        user.Age,
	}

	query, args, err := squirrel.Update(tableName).
		SetMap(setMap).
		Where(squirrel.Eq{
			primaryKey: strings.ToLower(user.Login),
		}).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumns)).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}

	queryer := i.getQueryer(nil)
	row := queryer.QueryRowxContext(ctx, query, args...)

	var (
		login, firstName, lastName string
		weight                     float32
		height, age                uint
	)

	err = row.Scan(&login, &firstName, &lastName, &weight, &height, &age)
	if err != nil {
		return models.User{}, errors.Wrap(errUpdate, err.Error())
	}
	updatedUser := models.User{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Weight:    weight,
		Height:    height,
		Age:       age}

	i.SetRedisData(login, updatedUser)
	return updatedUser, nil
}

func (i *implementation) Delete(ctx context.Context, login string) error {
	i.poolCh <- struct{}{}
	i.mu.Lock()
	defer func() {
		i.mu.Unlock()
		<-i.poolCh
	}()

	if _, err := i.getUserByLogin(login); err != nil {
		return errors.Wrap(errDelete, err.Error())
	}

	query, args, err := squirrel.Delete(tableName).
		Where(squirrel.Eq{
			primaryKey: strings.ToLower(login),
		}).
		Suffix(fmt.Sprintf("RETURNING %s", tableColumns)).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}

	queryer := i.getQueryer(nil)
	row := queryer.QueryRowxContext(ctx, query, args...)

	var (
		userLogin, firstName, lastName string
		weight                         float32
		height, age                    uint
	)

	err = row.Scan(&userLogin, &firstName, &lastName, &weight, &height, &age)
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}
	i.deleteRedisData(userLogin)

	return nil
}

func (i *implementation) List(_ context.Context, queryParams map[string]interface{}) (users []models.User, err error) {
	i.poolCh <- struct{}{}
	i.mu.RLock()
	defer func() {
		i.mu.RUnlock()
		<-i.poolCh
	}()

	selectBuilder := squirrel.Select(tableColumns).From(tableName)
	if val, ok := queryParams[QuerySortField]; ok {
		selectBuilder = selectBuilder.OrderBy(val.(string))
	}
	if val, ok := queryParams[QueryLimit]; ok {
		selectBuilder = selectBuilder.Limit(val.(uint64))
	}
	if val, ok := queryParams[QueryOffset]; ok {
		selectBuilder = selectBuilder.Offset(val.(uint64))
	}

	query, args, err := selectBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, errors.Wrap(errList, err.Error())
	}

	if err = i.db.Select(&users, query, args...); err != nil {
		return nil, errors.Wrap(errList, err.Error())
	}

	return users, nil
}

func (i *implementation) String(user models.User) string {
	return fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
		user.Login, user.FirstName, user.LastName, user.Weight, user.Height, user.Age)
}

func (i *implementation) getUserByLogin(login string) (models.User, error) {
	user, err := i.getRedisData(login)
	if err != nil {
		log.Info(err)
		return i.getUserByLoginFromDB(login)
	}
	return user, nil
}

func (i *implementation) getQueryer(tx *sqlx.Tx) sqlx.QueryerContext {
	if tx == nil {
		return i.db
	}
	return tx
}

func (i *implementation) getExecer(tx *sqlx.Tx) sqlx.ExecerContext {
	if tx == nil {
		return i.db
	}
	return tx
}

func (i *implementation) getUserByLoginFromDB(login string) (models.User, error) {
	query, args, err := squirrel.Select(tableColumns).
		From(tableName).
		Where(squirrel.Eq{
			primaryKey: strings.ToLower(login),
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return models.User{}, errors.Wrap(errSQL, err.Error())
	}

	var users []models.User
	if err = i.db.Select(&users, query, args...); err != nil {
		return models.User{}, errors.Wrap(errSQL, err.Error())
	}

	if users == nil {
		return models.User{}, errors.Wrapf(errUserNotExists, "user-login: [%s]", login)
	}
	if len(users) > 1 {
		return models.User{}, errors.Wrapf(errManyUsers, "user-login: [%s]", login)
	}
	i.SetRedisData(login, users[0])

	return users[0], nil
}

func (i *implementation) getRedisData(login string) (models.User, error) {
	getUser := i.redisClient.Get(login)
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

func (i *implementation) SetRedisData(key string, value interface{}) {
	status := i.redisClient.Set(key, value, redisExpiration)
	if status.Err() != nil {
		log.Info(status.Err())
	}
}

func (i *implementation) deleteRedisData(key string) {
	status := i.redisClient.Del(key)
	if status.Err() != nil {
		log.Info(status.Err())
	}
}
