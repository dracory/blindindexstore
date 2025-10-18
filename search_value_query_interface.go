package blindindexstore

import (
	"github.com/doug-martin/goqu/v9"
)

type SearchValueQueryInterface interface {
	ToSelectDataset(store StoreInterface) (*goqu.SelectDataset, []any, error)
	Validate() error

	ID() string
	HasID() bool
	SetID(id string) SearchValueQueryInterface

	IDIn() []string
	HasIDIn() bool
	SetIDIn(idIn []string) SearchValueQueryInterface

	SourceReferenceID() string
	HasSourceReferenceID() bool
	SetSourceReferenceID(sourceReferenceID string) SearchValueQueryInterface

	SearchValue() string
	HasSearchValue() bool
	SetSearchValue(searchValue string) SearchValueQueryInterface

	SearchType() string
	HasSearchType() bool
	SetSearchType(searchType string) SearchValueQueryInterface

	// Meta search methods
	AddMetaSearch(needle string) SearchValueQueryInterface
	GetMetaSearch() []string
	AddMetaSearchNot(needle string) SearchValueQueryInterface
	GetMetaSearchNot() []string

	Offset() int
	HasOffset() bool
	SetOffset(offset int) SearchValueQueryInterface

	Limit() int
	HasLimit() bool
	SetLimit(limit int) SearchValueQueryInterface

	OrderBy() string
	HasOrderBy() bool
	SetOrderBy(orderBy string) SearchValueQueryInterface

	OrderDirection() string
	HasOrderDirection() bool
	SetOrderDirection(orderByDirection string) SearchValueQueryInterface

	CountOnly() bool
	HasCountOnly() bool
	SetCountOnly(countOnly bool) SearchValueQueryInterface

	WithSoftDeleted() bool
	HasWithSoftDeleted() bool
	SetWithSoftDeleted(withSoftDeleted bool) SearchValueQueryInterface

	Columns() []any
	HasColumns() bool
	SetColumns(columns []any) SearchValueQueryInterface
}
