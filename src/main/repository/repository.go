package repository

import (
	"cabinet/src/main/datasource"
	"cabinet/src/main/model"
	"errors"

	"github.com/google/uuid"
)

type ProfileRepo struct {
	datasource *datasource.Datasource
}

func (p *ProfileRepo) FindById(uuid uuid.UUID) (*model.Profile, error) {
	if p.datasource == nil || p.datasource.Db == nil {
		return nil, errors.New("datasource is nil")
	}

	var profile = model.Profile{}

	err := p.datasource.Db.NewSelect().Model(&profile).Where("id = ?", uuid).Scan(p.datasource.Context)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}
