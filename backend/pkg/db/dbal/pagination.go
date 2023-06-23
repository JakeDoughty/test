package dbal

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	PageSize  int    `json:"pageSize"`
	PageIndex int    `json:"pageIndex"`
	Sort      string `json:"sort"`
}

func (p *Pagination) GetPage() int {
	if p.PageIndex == 0 {
		p.PageIndex = 1
	}
	return p.PageIndex
}
func (p *Pagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "ID desc"
	}
	return p.Sort
}

type List[T any] struct {
	Items      []*T        `json:"items,omitempty"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
	Pagination *Pagination `json:"pagination"`
}

func (list *List[T]) ScopePaginate(db *gorm.DB, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	var totalItems int64
	db.Model(list.Items).Count(&totalItems)

	list.Pagination = pagination
	list.TotalItems = totalItems
	list.TotalPages = int(math.Ceil(float64(totalItems) / float64(pagination.GetLimit())))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).
			Limit(pagination.GetLimit()).
			Order(pagination.GetSort())
	}
}
func (list *List[T]) ApplyPagination(db *gorm.DB, pagination *Pagination) *gorm.DB {
	return db.Scopes(list.ScopePaginate(db, pagination))
}
