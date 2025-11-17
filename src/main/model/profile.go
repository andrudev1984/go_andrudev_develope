package model

import (
	"cabinet/src/main/common"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Profile user data
type Profile struct {
	bun.BaseModel `bun:"table:users.profiles"`
	common.Modifiable
	Login        string         `bun:"type:varchar(50),notnull,unique"` // Login info
	FistName     string         `bun:"type:varchar(100),notnull"`
	MiddleName   string         `bun:"type:varchar(100),notnull"`
	LastName     string         `bun:"type:varchar(100),notnull"`
	Private      bool           `bun:"type:boolean,default:true"`
	PrimaryEmail string         `bun:"type:varchar(50),notnull,unique"`                     // Primary email, verified
	Email        []string       `bun:"type:varchar(50)[],array,default:array[]::varchar[]"` // Additional emails
	Phone        string         `bun:"type:varchar(50)"`
	Tags         []string       `bun:"type:varchar(50)[],array,default:array[]::varchar[]"`
	Biography    string         `bun:"type:text"`
	Company      string         `bun:"type:varchar(100)"`
	Location     string         `bun:"type:varchar(255)"`
	ExternalID   uuid.UUID      `bun:"type:uuid"`  // Keycloak id
	Avatar       uuid.UUID      `bun:"type:uuid"`  // S3 resource key
	Metadata     map[string]any `bun:"type:jsonb"` // Custom metadata
	Attachments  []*Attachment  `bun:"rel:has-many,join:id=user_id"`
}

// Attachment Profile custom material
type Attachment struct {
	bun.BaseModel `bun:"table:users.attachments"`
	common.NotModifiable
	common.Nameable
	Private  bool           `bun:"type:boolean,default:true"`
	Tags     []string       `bun:"type:varchar(50)[],array,default:array[]::varchar[]"`
	Title    string         `bun:"type:varchar(255),notnull"`
	S3Key    uuid.UUID      `bun:"type:uuid,notnull"`
	UserID   uuid.UUID      `bun:"type:uuid,notnull"`
	Metadata map[string]any `bun:"type:jsonb"` // Custom metadata
	Profile  Profile        `bun:"rel:belongs-to,join:user_id=id"`
}
