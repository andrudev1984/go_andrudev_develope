package test

import (
	"cabinet/src/main/model"
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stapelberg/postgrestest"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func TestTest(t *testing.T) {
	ctx := context.Background()
	srv, err := postgrestest.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}

	database, err := srv.NewDatabase(ctx)

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	bunDb := bun.NewDB(database, pgdialect.New())

	bunDb.RegisterModel((*model.Profile)(nil))

	bunDb.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	file, err := os.Open("../main/V1_Init.sql")
	migration, err := io.ReadAll(file)

	if err != nil {
		panic(err)
	} else {
		log.Println("Reading migrations... ok")
	}

	_, err = database.Exec(string(migration))
	if err != nil {
		panic(err)
	} else {
		log.Println("Loading migrations... ok")
	}

	fixture := dbfixture.New(bunDb, dbfixture.WithRecreateTables())
	err = fixture.Load(ctx, os.DirFS("./template"), "profiles.yaml")

	if err != nil {
		panic(err)
	} else {
		log.Println("Loading templates... ok")
	}

	var profiles []model.Profile

	err = bunDb.NewSelect().Column("id", "created", "login", "changed").Model(&profiles).Scan(ctx)

	if err != nil {
		panic(err)
	} else {
		log.Println("Select templates... ok")
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, profiles)
	assert.Equal(t, len(profiles), 2)

	var profile = prepareEntity()

	_, err = bunDb.NewInsert().Model(profile).Exec(ctx)

	assert.NoError(t, err)

	log.Println("Inserting new... ok")

	bunDb.NewSelect().Column("id", "created", "login", "changed").Model(&profiles).Scan(ctx)

	assert.Equal(t, len(profiles), 3)
	assert.NotNil(t, profiles[2].ID)
	assert.NotNil(t, profiles[2].Created)
	assert.NotNil(t, profiles[2].Changed)

	profile.Metadata = map[string]interface{}{}
	profile.Metadata["aaa"] = "bbb"
	profile.Metadata["bbb"] = "ccc"
	profile.Metadata["ddd"] = "eee"

	_, err = bunDb.NewUpdate().Model(profile).Where("id = ?", profiles[2].ID).Exec(ctx)

	if err != nil {
		panic(err)
	} else {
		log.Println("Update templates... ok")
	}

	t.Cleanup(srv.Cleanup)
}

func prepareEntity() *model.Profile {
	var profile = &model.Profile{}

	profile.Login = "New login"
	profile.Location = "New location"
	profile.Biography = "New bio"
	profile.Company = "New company"
	profile.FistName = "New fist name"
	profile.MiddleName = "New middle name"
	profile.LastName = "New last name"
	profile.Avatar = uuid.New()
	profile.ExternalID = uuid.New()

	return profile
}
