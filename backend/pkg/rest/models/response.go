package models

import (
	"fmt"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/dbal"
	"github.com/gin-gonic/gin"
)

type Response interface {
	IsSucceeded() bool
	Send(c *gin.Context, status int)
}

type ErrorResponse struct {
	Succeeded bool   `json:"succeeded"`
	Message   string `json:"msg"`
}

func (resp *ErrorResponse) IsSucceeded() bool { return resp.Succeeded }
func (resp *ErrorResponse) Send(c *gin.Context, status int) {
	c.AbortWithStatusJSON(status, resp)
}

func NewErrorResponsef(format string, args ...any) *ErrorResponse {
	return &ErrorResponse{Message: fmt.Sprintf(format, args...)}
}
func NewErrorResponse(msg string) *ErrorResponse { return &ErrorResponse{Message: msg} }
func ErrorToResponse(err error) *ErrorResponse   { return NewErrorResponse(err.Error()) }

type SuccessResponse struct {
	Succeeded bool `json:"succeeded"`
	Data      any  `json:"data"`
}

func (resp *SuccessResponse) IsSucceeded() bool               { return true }
func (resp *SuccessResponse) Send(c *gin.Context, status int) { c.JSON(status, resp) }

func NewSuccessResponse(data any) *SuccessResponse {
	return &SuccessResponse{Succeeded: true, Data: data}
}

type List[T any] struct {
	Items      []*T  `json:"items"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
	PageIndex  int   `json:"pageIndex"`
	PageSize   int   `json:"pageSize"`
}
type ListResponse[T any] struct {
	Succeeded bool    `json:"succeeded"`
	Data      List[T] `json:"data"`
}

func NewListResponse[T any](list *dbal.List[T]) *ListResponse[T] {
	response := &ListResponse[T]{
		Succeeded: true,
		Data: List[T]{
			Items:      list.Items,
			TotalItems: list.TotalItems,
			TotalPages: list.TotalPages,
		},
	}
	if list.Pagination != nil {
		response.Data.PageIndex = list.Pagination.PageIndex
		response.Data.PageSize = list.Pagination.PageSize
	}
	return response
}
func ConvertToListResponse[T any, Entity any](list *dbal.List[Entity], convert func(*Entity) *T) *ListResponse[T] {
	response := &ListResponse[T]{
		Succeeded: true,
		Data: List[T]{
			TotalItems: list.TotalItems,
			TotalPages: list.TotalPages,
			Items:      make([]*T, len(list.Items)),
		},
	}

	if list.Pagination != nil {
		response.Data.PageIndex = list.Pagination.PageIndex
		response.Data.PageSize = list.Pagination.PageSize
	}

	for i, entity := range list.Items {
		response.Data.Items[i] = convert(entity)
	}
	return response
}
func (resp *ListResponse[T]) Send(c *gin.Context, status int) { c.JSON(status, resp) }
