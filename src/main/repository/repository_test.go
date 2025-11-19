package repository

import (
	"cabinet/src/main/datasource"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var dataSource *datasource.Datasource
var profileTestId1 = uuid.New()
var profileBio = "Profile bio"
var profileTestId2 = uuid.New()
var testMock sqlmock.Sqlmock

var findReqFormat = "SELECT \"profile\".\"id\", \"profile\".\"created\"," +
	" \"profile\".\"changed\"," +
	" \"profile\".\"login\"," +
	" \"profile\".\"fist_name\"," +
	" \"profile\".\"middle_name\"," +
	" \"profile\".\"last_name\"," +
	" \"profile\".\"private\"," +
	" \"profile\".\"primary_email\"," +
	" \"profile\".\"email\"," +
	" \"profile\".\"phone\"," +
	" \"profile\".\"tags\"," +
	" \"profile\".\"biography\"," +
	" \"profile\".\"company\"," +
	" \"profile\".\"location\"," +
	" \"profile\".\"external_id\"," +
	" \"profile\".\"avatar\"," +
	" \"profile\".\"metadata\" FROM \"users\".\"profiles\" AS \"profile\" WHERE (id = '%s')"

func TestMain(m *testing.M) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	testMock = mock

	if err != nil {
		slog.Error("An error was not expected when opening a stub database connection",
			slog.Any("err", err.Error()))
		panic(err)
	}

	defer db.Close()

	var bunDb = bun.NewDB(db, pgdialect.New())

	bunDb.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	dataSource = &datasource.Datasource{Db: bunDb, Context: context.Background()}

	m.Run()
}

func TestFindProfileById(t *testing.T) {
	var rows = testMock.NewRows([]string{"id", "biography"})
	rows.AddRow(profileTestId1, profileBio)

	testMock.ExpectQuery(fmt.Sprintf(findReqFormat, profileTestId1)).WillReturnRows(rows)
	testMock.ExpectQuery(fmt.Sprintf(findReqFormat, profileTestId2)).WillReturnRows(testMock.NewRows([]string{"id"}))

	var profileRepo = &ProfileRepo{datasource: nil}

	var profile, err = profileRepo.FindById(profileTestId1)
	assert.Nil(t, profile)
	assert.Error(t, err)
	slog.Error("An error was not expected when finding profile", slog.Any("err", err.Error()))

	profileRepo = &ProfileRepo{datasource: dataSource}

	profile, err = profileRepo.FindById(profileTestId1)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, profileTestId1, profile.ID)
	assert.Equal(t, profileBio, profile.Biography)

	profile, err = profileRepo.FindById(profileTestId2)
	assert.Nil(t, profile)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	slog.Error("An error was not expected when finding profile", slog.Any("err", err.Error()))

	slog.Info("TestFindProfileById is successful")
}
