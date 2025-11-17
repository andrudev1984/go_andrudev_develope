package test

import (
	"cabinet/src/main/model"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
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

const MigrationsPath = "../main/migrations"

func TestM(t *testing.T) {
	ctx := context.Background()
	pgt, err := postgrestest.Start(ctx)
	if err != nil {
		panic(err)
	}

	sqlDb, err := pgt.NewDatabase(ctx)

	err = sqlDb.Ping()
	if err != nil {
		panic(err)
	}

	bunDb := bun.NewDB(sqlDb, pgdialect.New())

	bunDb.RegisterModel((*model.Profile)(nil))
	bunDb.RegisterModel((*model.Attachment)(nil))

	bunDb.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	defer pgt.Cleanup()

	t.Run("Migrations", func(t *testing.T) {
		dir, err2 := os.ReadDir(MigrationsPath)
		assert.NoError(t, err2)

		for _, d := range dir {
			file, err := os.Open(filepath.Join(MigrationsPath, d.Name()))
			migration, err := io.ReadAll(file)

			log.Println(string(migration))

			if err != nil {
				panic(err)
			} else {
				log.Printf("Reading migration... %s ok\n", d.Name())
			}

			_, err = sqlDb.Exec(string(migration))
			if err != nil {
				panic(err)
			} else {
				log.Printf("Loading migration... %s ok\n", d.Name())
			}
		}

		fixture := dbfixture.New(bunDb, dbfixture.WithTruncateTables())
		err = fixture.Load(ctx, os.DirFS("./template"), "profiles.yaml")

		if err != nil {
			panic(err)
		} else {
			log.Println("Loading templates... ok")
		}
	})

	t.Run("Get profiles", func(t *testing.T) {
		var profiles []model.Profile

		err = bunDb.NewSelect().Model(&profiles).Scan(ctx)

		if err != nil {
			panic(err)
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, profiles)
		assert.Equal(t, len(profiles), 2)

		log.Println("Select profiles... ok")
	})

	t.Run("Get attachments", func(t *testing.T) {
		var aAttachments []model.Attachment

		err = bunDb.NewSelect().Model(&aAttachments).Scan(ctx)

		if err != nil {
			panic(err)
		}

		assert.Equal(t, 2, len(aAttachments))

		log.Println("Select attachments... ok")
	})

	t.Run("Operate profile", func(t *testing.T) {
		var profiles []model.Profile

		var profile = prepareProfileEntity()

		_, err = bunDb.NewInsert().Model(profile).Exec(ctx)

		assert.NoError(t, err)

		err = bunDb.NewSelect().Model(&profiles).Scan(ctx)

		if err != nil {
			panic(err)
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, profiles)
		assert.Equal(t, len(profiles), 3)

		profile.Metadata = map[string]interface{}{}
		profile.Metadata["aaa"] = "bbb"
		profile.Metadata["bbb"] = "ccc"
		profile.Metadata["ddd"] = "eee"
		profile.Tags = []string{"tag1", "tag2"}

		_, err = bunDb.NewUpdate().Model(profile).Where("id = ?", profiles[2].ID).Exec(ctx)

		if err != nil {
			panic(err)
		}

		_, err = bunDb.NewDelete().Model(profile).Where("id = ?", profiles[2].ID).Exec(ctx)

		if err != nil {
			panic(err)
		}

		err = bunDb.NewSelect().Model(&profiles).Scan(ctx)

		if err != nil {
			panic(err)
		}

		assert.Equal(t, len(profiles), 2)

		log.Println("Operating new profile... ok")
	})

	t.Run("Operate attachments", func(t *testing.T) {
		var profiles []model.Profile

		bunDb.NewSelect().Model(&profiles).Relation("Attachments").Scan(ctx)

		assert.Equal(t, 1, len(profiles[0].Attachments))
		assert.Equal(t, 1, len(profiles[1].Attachments))

		var attachment = prepareAttachmentEntity(uuid.New())

		_, err = bunDb.NewInsert().Model(attachment).Exec(ctx)

		assert.Error(t, err)

		log.Println(err)

		attachment.UserID = profiles[0].ID

		_, err := bunDb.NewInsert().Model(attachment).Exec(ctx)

		assert.NoError(t, err)

		count, _ := bunDb.NewSelect().Model(attachment).Count(ctx)

		assert.Equal(t, 3, count)

		bunDb.NewSelect().Model(attachment).Relation("Profile").Where("s3_key = ?", attachment.S3Key).Scan(ctx)

		assert.Equal(t, attachment.UserID, attachment.Profile.ID)

		attachment.Metadata = map[string]interface{}{}

		attachment.Metadata["aaa"] = "bbb"
		attachment.Metadata["bbb"] = "ccc"
		attachment.Metadata["ddd"] = "eee"
		attachment.Tags = []string{"tag1", "tag2"}

		_, err = bunDb.NewUpdate().Model(attachment).Where("id = ?", attachment.ID).Exec(ctx)

		assert.NoError(t, err)

		_, err = bunDb.NewDelete().Model(attachment).Where("id = ?", attachment.ID).Exec(ctx)

		assert.NoError(t, err)

		log.Println("Operating new attachment... ok")

	})

	pgt.Cleanup()
}

func prepareProfileEntity() *model.Profile {
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

func prepareAttachmentEntity(userId uuid.UUID) *model.Attachment {
	var attachment = &model.Attachment{}

	attachment.Name = "New filename"
	attachment.Description = "New description"
	attachment.S3Key = uuid.New()
	attachment.Title = "New title"
	attachment.UserID = userId
	attachment.Private = false

	return attachment
}
