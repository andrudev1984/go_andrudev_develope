package common

import (
	"cabinet/src/main/model"
	"cabinet/src/main/model/common"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestIdInfo(test *testing.T) {
	var profile = &model.Profile{}
	profile.ID = uuid.New()

	var attachment = &model.Attachment{}
	attachment.ID = uuid.New()

	var idInfo1 = &IdInfo{}

	idInfo1.From(nil)

	assert.Empty(test, idInfo1.ID)

	idInfo1.From(profile)

	assert.NotEmpty(test, idInfo1.ID)
	assert.Equal(test, profile.ID, idInfo1.ID)

	var idInfo2 = &IdInfo{}
	idInfo2.From(attachment)

	assert.NotEmpty(test, idInfo2.ID)
	assert.Equal(test, attachment.ID, idInfo2.ID)

	slog.Info("TestIdInfo success")
}

func TestNameInfo(test *testing.T) {
	var profile = &model.Profile{}
	profile.ID = uuid.New()
	profile.FistName = "Fist"
	profile.MiddleName = "Middle"
	profile.LastName = "Last"

	var nameInfo = &ShortNamedInfo{}

	nameInfo.From(nil)

	assert.Empty(test, nameInfo.IdInfo.ID)
	assert.Empty(test, nameInfo.Name)
	assert.Empty(test, nameInfo.Description)

	nameInfo.From(profile)

	assert.NotEmpty(test, nameInfo.IdInfo.ID)
	assert.Equal(test, profile.ID, nameInfo.IdInfo.ID)
	assert.Equal(test, profile.FullName(), nameInfo.Name)
	assert.Empty(test, nameInfo.Description)

	var attachment = &model.Attachment{}
	attachment.ID = uuid.New()
	attachment.Name = "New Attachment"
	attachment.Description = "New Attachment Description"

	var nameInfo1 = &ShortNamedInfo{}
	nameInfo1.From(attachment)

	assert.NotEmpty(test, nameInfo1.IdInfo.ID)
	assert.Equal(test, attachment.ID, nameInfo1.IdInfo.ID)
	assert.Equal(test, attachment.Name, nameInfo1.Name)
	assert.Equal(test, attachment.Description, nameInfo1.Description)

	slog.Info("TestNameInfo success")
}

func TestFrom(test *testing.T) {
	var profile = &model.Profile{}
	profile.ID = uuid.New()
	profile.FistName = "Fist"
	profile.MiddleName = "Middle"
	profile.LastName = "Last"

	var nameInfo = &ShortNamedInfo{}

	nameInfo.From(profile)

	var result = ResultDto[ShortNamedInfo]{Result: *nameInfo}

	assert.NotEmpty(test, result.Result)
	assert.Equal(test, nameInfo.Name, result.Result.Name)
	assert.Equal(test, nameInfo.ID, result.Result.ID)

	content, err := json.Marshal(result)

	assert.NoError(test, err)
	slog.Info("Json result content", slog.String("content", string(content)))

	slog.Info("TestResult success")
}

func TestError(test *testing.T) {
	var errorResult = &ErrorDto{}

	errorResult.From(401, "Not Authorized")

	assert.NotEmpty(test, errorResult.Code)
	assert.NotEmpty(test, errorResult.Message)
	assert.Empty(test, errorResult.Details)

	errorResult.From(404, "Not Found", "File is not found")
	assert.NotEmpty(test, errorResult.Details)

	errorResult = BuildError(404, "Not Found", "File is not found")

	assert.Equal(test, uint16(404), errorResult.Code)
	assert.Equal(test, "Not Found", errorResult.Message)
	assert.Equal(test, []string{"File is not found"}, errorResult.Details)

	slog.Info("TestError success")
}

func TestPagination(test *testing.T) {
	var pagination = BuildPagination(10, 100, 10)

	assert.NotNil(test, pagination)
	assert.Equal(test, uint(10), pagination.Page)
	assert.Equal(test, uint64(100), pagination.Total)
	assert.Equal(test, uint(10), pagination.PageSize)

	assert.False(test, IsValidPagination(&Pagination{
		Page:     10,
		Total:    10,
		PageSize: 10,
	}))
	assert.Nil(test, BuildPagination(10, 0, 10))

	slog.Info("TestPagination success")
}

func TestBuildPaged(test *testing.T) {
	var pagination = &Pagination{10, 0, 10}

	var paged = BuildPaged[ShortNamedInfo](nil, pagination)

	assert.Nil(test, paged)

	paged = BuildPaged[ShortNamedInfo](nil, BuildPagination(0, 1, 1))

	assert.NotNil(test, paged)
	assert.Equal(test, 0, len(paged.ResultDto.Entities))

	var indents []common.Identifiable

	indents = append(indents, common.Identifiable{ID: uuid.New()})
	indents = append(indents, common.Identifiable{ID: uuid.New()})

	var paged1 = BuildPaged[common.Identifiable](indents, BuildPagination(1, 2, 1))

	assert.NotNil(test, paged1)
	assert.Equal(test, 1, len(paged1.ResultDto.Entities))
	assert.Equal(test, indents[1].ID, paged1.ResultDto.Entities[0].ID)

	content, err := json.Marshal(paged1)

	assert.NoError(test, err)
	slog.Info("Json result content", slog.String("content", string(content)))

	slog.Info("TestBuildPaged success")
}
