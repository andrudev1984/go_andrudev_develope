package model

import (
	"cabinet/src/main/common"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// Profile user data
type Profile struct {
	bun.BaseModel `bun:"table:users.profiles"`
	common.Modifiable
	Login        string         `bun:"type:varchar(50),notnull,scanonly,unique"` // Login info
	FistName     string         `bun:"type:varchar(100),notnull,default:''"`
	MiddleName   string         `bun:"type:varchar(100),notnull,default:''"`
	LastName     string         `bun:"type:varchar(100),notnull,default:''"`
	PrimaryEmail string         `bun:"type:email,notnull,unique"` // Primary email, verified
	Email        []string       `bun:"type:email,array"`          // Additional emails
	Phone        string         `bun:"type:phone"`
	Tags         []string       `bun:"type:tags,array"`
	Biography    string         `bun:"type:text,default:''"`
	Company      string         `bun:"type:varchar(100),default:''"`
	Location     string         `bun:"type:varchar(255),default:''"`
	ExternalID   uuid.UUID      `bun:"type:uuid"`                        // Keycloak id
	Avatar       uuid.UUID      `bun:"type:uuid"`                        // S3 resource key
	Metadata     map[string]any `bun:"type:metadata,jsonb,default:'{}'"` // Custom metadata
	Attachments  []*Attachment  `bun:"rel:has-many,join:id=user_id"`
}

// Attachment Profile custom material
type Attachment struct {
	common.NotModifiable
	common.Nameable
	Private bool      `bun:"type:boolean,default:true"`
	Tags    []string  `bun:"type:tags,array"`
	Title   string    `bun:"title,notnull"`
	S3Key   uuid.UUID `bun:"type:uuid, notnull"`
	UserID  uuid.UUID `bun:"type:uuid,notnull"`
	Profile Profile   `bun:"rel:belongs-to,join:user_id=id"`
}

// Created / Changed hooks
var _ bun.BeforeAppendModelHook = (*Profile)(nil)
var _ bun.BeforeAppendModelHook = (*Attachment)(nil)

func (p Profile) BeforeAppendModel(ctx context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		p.Created = time.Now().UTC()
		p.Changed = time.Now().UTC()
	case *bun.UpdateQuery:
		p.Changed = time.Now().UTC()
	}
	return nil
}

func (a Attachment) BeforeAppendModel(ctx context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		a.Created = time.Now().UTC()
	}
	return nil
}
