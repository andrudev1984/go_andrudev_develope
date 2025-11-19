package main

import (
	"cabinet/src/main/model"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func main() {
	p := &model.Profile{
		Login:        "",
		FistName:     "",
		MiddleName:   "",
		LastName:     "",
		PrimaryEmail: "",
		Email:        nil,
		Phone:        "",
		Tags:         nil,
		Biography:    "",
		Company:      "",
		Location:     "",
		ExternalID:   uuid.UUID{},
		Avatar:       uuid.UUID{},
		Metadata:     nil,
		Attachments:  nil,
	}

	p.ID = uuid.New()
	p.Created = time.Now().UTC()
	p.Changed = time.Now().UTC()

	a := &model.Attachment{
		Private: false,
		Tags:    nil,
		Title:   "",
		Profile: model.Profile{},
	}

	a.ID = uuid.New()
	a.Created = time.Now().UTC()
	a.Name = "Name"
	a.Description = "Description"
	a.S3Key = uuid.New()
	a.UserID = uuid.New()

	pJson, _ := json.Marshal(p)
	aJson, _ := json.Marshal(a)

	slog.Info(string(pJson))
	slog.Info(string(aJson))
}
