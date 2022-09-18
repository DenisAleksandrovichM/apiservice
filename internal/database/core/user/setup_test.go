package user

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"testing"
)

type usersRepoFixtures struct {
	usersRepo User
	db        *sqlx.DB
	dbMock    sqlmock.Sqlmock
}

func setUp(t *testing.T) usersRepoFixtures {
	var fixture usersRepoFixtures

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	fixture.db = sqlx.NewDb(db, "sqlmock")
	fixture.dbMock = mock
	fixture.usersRepo = New(fixture.db)
	return fixture
}

func (f *usersRepoFixtures) tearDown() {
	f.db.Close()
}
