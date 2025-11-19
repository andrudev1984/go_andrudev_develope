package common

import (
	"cabinet/src/main/model/interfaces"
	"log/slog"

	"github.com/google/uuid"
)

type IdInfo struct {
	ID uuid.UUID `json:"id"`
}

func (i *IdInfo) From(identifiable interfaces.Identifiable) {
	if identifiable == nil {
		slog.Error("Nil identifiable")
		return
	}

	i.ID = identifiable.GetId()
}

type ShortNamedInfo struct {
	IdInfo
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (i *ShortNamedInfo) From(nameable interfaces.Nameable) {
	if nameable == nil {
		return
	}
	i.Name = nameable.GetName()
	i.Description = nameable.GetDescription()
	i.IdInfo.From(nameable)
}

type ResultDto[T any] struct {
	Result T `json:"result"` // result object
}

type ErrorDto struct {
	Code    uint16   `json:"code"`    // error code
	Message string   `json:"message"` // error message
	Details []string `json:"details"` // error details
}

func (e *ErrorDto) From(code uint16, message string, details ...string) {
	e.Code = code
	e.Message = message
	e.Details = details
}

func BuildError(code uint16, message string, details ...string) *ErrorDto {
	return &ErrorDto{code, message, details}
}

type Pagination struct {
	Page     uint   `json:"page"`
	Total    uint64 `json:"total"`
	PageSize uint   `json:"pageSize"`
}

func IsValidPagination(pagination *Pagination) bool {
	return pagination != nil && (uint64)((pagination.Page+1)*pagination.PageSize) <= pagination.Total
}

func BuildPagination(page uint, total uint64, pageSize uint) *Pagination {
	if (uint64)(page*pageSize) > total {
		slog.Error("Incorrect Pagination params",
			slog.Any("page", page),
			slog.Any("total", total),
			slog.Any("pageSize", pageSize))
		return nil
	}
	return &Pagination{page, total, pageSize}
}

type Paged[T any] struct {
	Entities []T        `json:"entities"`
	Pageable Pagination `json:"pageable"`
}

type PagedResult[T any] struct {
	ResultDto Paged[T] `json:"result"`
}

func BuildPaged[T any](entities []T, pageable *Pagination) *PagedResult[T] {
	if !IsValidPagination(pageable) {
		if pageable != nil {
			slog.Error("Incorrect Pagination params",
				slog.Any("page", pageable.Page),
				slog.Any("total", pageable.Total),
				slog.Any("pageSize", pageable.PageSize))
		}

		return nil
	}

	var paged = &Paged[T]{}

	if entities == nil || len(entities) == 0 {
		paged.Entities = []T{}
		paged.Pageable = *pageable

		return &PagedResult[T]{*paged}
	}

	paged.Entities = entities[pageable.PageSize*(pageable.Page) : pageable.PageSize*(pageable.Page+1)]
	paged.Pageable = *pageable

	return &PagedResult[T]{*paged}
}
