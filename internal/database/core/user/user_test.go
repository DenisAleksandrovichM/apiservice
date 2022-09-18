package user

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DenisAleksandrovichM/apiservice/pkg/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

const selectQuery = `SELECT login, first_name, last_name, weight, height, age FROM users WHERE login = $1`

var (
	errSome = errors.New("some error")
	user1   = models.User{
		Login:     "test_login",
		FirstName: "test_fn",
		LastName:  "test_ln",
		Weight:    80,
		Height:    180,
		Age:       60,
	}
	user2 = models.User{
		Login:     "test_login1",
		FirstName: "test_fn1",
		LastName:  "test_ln1",
		Weight:    81,
		Height:    181,
		Age:       61,
	}
)

func TestNewUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{})
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnRows(rows)
		queryStore = regexp.QuoteMeta(`INSERT INTO users (login, first_name, last_name, weight, height, age) VALUES ($1,$2,$3,$4,$5,$6)`)
		rows = sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age).
			WillReturnRows(rows)
		// act
		result, err := f.usersRepo.Create(context.Background(), user1)

		// assert
		require.NoError(t, err)
		assert.Equal(t, user1, result)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{})
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnRows(rows)
		queryStore = regexp.QuoteMeta(`INSERT INTO users (login, first_name, last_name, weight, height, age) VALUES ($1,$2,$3,$4,$5,$6)`)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age).
			WillReturnError(errSome)
		// act
		_, err := f.usersRepo.Create(context.Background(), user1)

		// assert
		assert.EqualError(t, err, errors.Wrap(errAdd, errSome.Error()).Error())
	})
}

func TestReadUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnRows(rows)

		// act
		result, err := f.usersRepo.Read(context.Background(), user1.Login)

		// assert
		require.NoError(t, err)
		assert.Equal(t, user1, result)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnError(errSome)

		// act
		_, err := f.usersRepo.Read(context.Background(), user1.Login)

		// assert
		assert.EqualError(t, err, errors.Wrap(errors.Wrap(errRead, errSQL.Error()), errSome.Error()).Error())
	})
}

func TestUpdateUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs("test_login").
			WillReturnRows(rows)
		queryStore = regexp.QuoteMeta(`UPDATE users SET age = $1, first_name = $2, height = $3, last_name = $4, weight = $5 WHERE login = $6`)
		rows = sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Age, user1.FirstName, user1.Height, user1.LastName, user1.Weight, user1.Login).
			WillReturnRows(rows)
		// act
		result, err := f.usersRepo.Update(context.Background(), user1)

		// assert
		require.NoError(t, err)
		assert.Equal(t, user1, result)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{})
		f.dbMock.ExpectQuery(queryStore).
			WithArgs("test_login").
			WillReturnRows(rows)
		queryStore = regexp.QuoteMeta(`UPDATE users SET age = $1, first_name = $2, height = $3, last_name = $4, weight = $5 WHERE login = $6`)
		rows = sqlmock.NewRows([]string{"age", "first_name", "height", "last_name", "weight", "login"}).
			AddRow(user1.Age, user1.FirstName, user1.Height, user1.LastName, user1.Weight, user1.Login)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Age, user1.FirstName, user1.Height, user1.LastName, user1.Weight, user1.Login).
			WillReturnError(errSome)
		// act
		_, err := f.usersRepo.Update(context.Background(), user1)

		// assert
		assert.ErrorIs(t, err, errUpdate)
	})
}

func TestDeleteUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnRows(rows)

		queryStore = regexp.QuoteMeta(`DELETE FROM users WHERE login = $1`)
		rows = sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnRows(rows)

		// act
		err := f.usersRepo.Delete(context.Background(), user1.Login)

		// assert
		require.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(selectQuery)
		rows := sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnRows(rows)

		queryStore = regexp.QuoteMeta(`DELETE FROM users WHERE login = $1`)
		rows = sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		f.dbMock.ExpectQuery(queryStore).
			WithArgs(user1.Login).
			WillReturnError(errSome)

		// act
		err := f.usersRepo.Delete(context.Background(), user1.Login)

		// assert
		assert.EqualError(t, err, errors.Wrap(errDelete, errSome.Error()).Error())
	})
}

func TestListUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(`SELECT login, first_name, last_name, weight, height, age FROM users`)
		rows := sqlmock.NewRows([]string{"login", "first_name", "last_name", "weight", "height", "age"}).
			AddRow(user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age).
			AddRow(user2.Login, user2.FirstName, user2.LastName, user2.Weight, user2.Height, user2.Age)
		f.dbMock.ExpectQuery(queryStore).
			WillReturnRows(rows)

		// act
		result, err := f.usersRepo.List(context.Background(), nil)

		// assert
		require.NoError(t, err)
		assert.Equal(t, []models.User{user1, user2}, result)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryStore := regexp.QuoteMeta(`SELECT login, first_name, last_name, weight, height, age FROM users`)
		f.dbMock.ExpectQuery(queryStore).
			WillReturnError(errSome)
		// act
		_, err := f.usersRepo.List(context.Background(), nil)

		// assert
		assert.EqualError(t, err, errors.Wrap(errList, errSome.Error()).Error())
	})

}

func TestStringUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		// act
		result := f.usersRepo.String(user1)

		// assert
		expectedStr := fmt.Sprintf("Login: %s, first name: %s, last name: %s,\nweight: %.2f, height: %d, age: %d",
			user1.Login, user1.FirstName, user1.LastName, user1.Weight, user1.Height, user1.Age)
		assert.Equal(t, expectedStr, result)
	})
}
