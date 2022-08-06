package local

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"sync"

	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/cache"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/core/user/models"
)

const (
	poolSize       = 10
	tableName      = "users"
	tableColumns   = "login, first_name, last_name, weight, height, age"
	primaryKey     = "login"
	QuerySortField = "SortField"
	QueryLimit     = "Limit"
	QueryOffset    = "Offset"
)

var (
	errUserNotExists = errors.New("user does not exist")
	errUserExists    = errors.New("user exists")
	errManyUsers     = errors.New("there are many users with this login")
	errSQL           = errors.New("SQL error")
	errAdd           = errors.New("add error")
	errRead          = errors.New("read error")
	errUpdate        = errors.New("update error")
	errDelete        = errors.New("delete error")
	errList          = errors.New("list error")
)

func New(pool *pgxpool.Pool) cachePkg.Interface {
	return &cache{
		mu:     sync.RWMutex{},
		data:   map[string]models.User{},
		poolCh: make(chan struct{}, poolSize),
		pool:   pool,
	}
}

type cache struct {
	mu     sync.RWMutex
	data   map[string]models.User
	poolCh chan struct{}
	pool   *pgxpool.Pool
}

func (c *cache) List(ctx context.Context, queryParams map[string]interface{}) (users []models.User, err error) {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		<-c.poolCh
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

	if err = pgxscan.Select(ctx, c.pool, &users, query, args...); err != nil {
		return nil, errors.Wrap(errList, err.Error())
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

	_, err := c.getUserByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, errUserNotExists) {
		return errors.Wrap(errAdd, err.Error())
	} else if err == nil {
		return errors.Wrap(errAdd, errUserExists.Error())
	}

	query, args, err := squirrel.Insert(tableName).
		Columns(tableColumns).
		Values(strings.ToLower(user.Login), user.FirstName, user.LastName, user.Weight, user.Height, user.Age).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(errAdd, err.Error())
	}

	if _, err = c.pool.Exec(ctx, query, args...); err != nil {
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

	user, err := c.getUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, errors.Wrap(errRead, err.Error())
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

	_, err := c.getUserByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, errUserNotExists) {
		return errors.Wrap(errAdd, err.Error())
	} else if err == nil {
		return errors.Wrap(errAdd, errUserExists.Error())
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
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(errUpdate, err.Error())
	}

	if _, err = c.pool.Exec(ctx, query, args...); err != nil {
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

	if _, err := c.getUserByLogin(ctx, login); err != nil {
		return errors.Wrap(errDelete, err.Error())
	}

	query, args, err := squirrel.Delete(tableName).
		Where(squirrel.Eq{
			primaryKey: strings.ToLower(login),
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(errDelete, err.Error())
	}

	if _, err = c.pool.Exec(ctx, query, args...); err != nil {
		return errors.Wrap(errDelete, err.Error())
	}

	return nil
}

func (c *cache) getUserByLogin(ctx context.Context, login string) (models.User, error) {
	query, args, err := squirrel.Select(tableColumns).
		From(tableName).
		Where(squirrel.Eq{
			primaryKey: strings.ToLower(login),
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return models.User{}, errors.Wrap(errSQL, err.Error())
	}

	var users []models.User
	if err = pgxscan.Select(ctx, c.pool, &users, query, args...); err != nil {
		return models.User{}, errors.Wrap(errSQL, err.Error())
	}

	if users == nil {
		return models.User{}, errors.Wrapf(errUserNotExists, "user-login: [%s]", login)
	}
	if len(users) > 1 {
		return models.User{}, errors.Wrapf(errManyUsers, "user-login: [%s]", login)
	}

	return users[0], nil
}
