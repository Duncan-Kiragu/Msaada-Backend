package filter

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var orders = []string{"asc", "desc"}

func NewFilter() *Filter {
	return &Filter{
		Search: "",
		Page:   0,
		Limit:  0,
		Sort:   os.Getenv("API_DEFAULT_SORT"),
		Order:  os.Getenv("API_DEFAULT_ORDER"),
	}
}

type (
	Filter struct {
		Search string `query:"search" form:"search" example:"name"`
		Page   int    `query:"page" form:"page" example:"1"`
		Limit  int    `query:"limit" form:"limit" example:"10"`
		Sort   string `query:"sort" form:"sort" example:"'updated_at', 'created_at', 'name' or some other field of the response object"`
		Order  string `query:"order" form:"order" example:"descending order 'desc' or ascending order 'asc'"`
	}

	UserFilter struct {
		Filter
		ProfileID uint `query:"profile_id" form:"profile_id" example:"1"`
	}
)

func (s *Filter) ApplySearchLike(db *gorm.DB, columns ...string) *gorm.DB {
	if len(columns) > 0 && s.Search != "" {
		whereLike := func(column, value string) string {
			return fmt.Sprintf("unaccent(LOWER(%v)) LIKE unaccent(LOWER('%%%v%%'))", column, value)
		}
		where := ""

		for i, column := range columns {
			if i > 0 {
				where = where + " or "
			}

			where = where + whereLike(column, s.Search)
		}

		if where != "" {
			return db.Where(where)
		}
	}

	return db
}

func (s *Filter) ApplyOrder(db *gorm.DB) *gorm.DB {
	s.check()
	return db.Order(fmt.Sprintf("%v %v", s.Sort, s.Order))
}

func (s *Filter) ApplyPagination(db *gorm.DB) *gorm.DB {
	if s.Page > 0 && s.Limit > 0 {
		return db.Offset((s.Page - 1) * s.Limit).Limit(s.Limit)
	}

	return db
}

func (s *Filter) check() {
	s.Order = strings.ToLower(s.Order)
	if !slices.Contains(orders, s.Order) {
		s.Order = os.Getenv("API_DEFAULT_ORDER")
	}
}
