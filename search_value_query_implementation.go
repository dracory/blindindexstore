package blindindexstore

import (
	"errors"
	"strings"
)

// == INTERFACE ===============================================================

type SearchValueQueryInterface interface {
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

// == TYPE ====================================================================

func SearchValueQuery() SearchValueQueryInterface {
	return NewSearchValueQuery()
}

func NewSearchValueQuery() SearchValueQueryInterface {
	return &searchValueQueryImplementation{}
}

type searchValueQueryImplementation struct {
	id    string
	hasID bool

	idIn    []string
	hasIDIn bool

	sourceReferenceID    string
	hasSourceReferenceID bool

	searchValue    string
	hasSearchValue bool

	searchType    string
	hasSearchType bool

	offset    int
	hasOffset bool

	limit    int
	hasLimit bool

	orderBy    string
	hasOrderBy bool

	orderDirection    string
	hasOrderDirection bool

	countOnly    bool
	hasCountOnly bool

	withSoftDeleted    bool
	hasWithSoftDeleted bool

	columns    []any
	hasColumns bool
}

func (o *searchValueQueryImplementation) ID() string {
	return o.id
}

func (o *searchValueQueryImplementation) HasID() bool {
	return o.hasID
}

func (o *searchValueQueryImplementation) SetID(id string) SearchValueQueryInterface {
	o.id = id
	o.hasID = true
	return o
}

func (o *searchValueQueryImplementation) IDIn() []string {
	return o.idIn
}

func (o *searchValueQueryImplementation) HasIDIn() bool {
	return o.hasIDIn
}

func (o *searchValueQueryImplementation) SetIDIn(idIn []string) SearchValueQueryInterface {
	o.idIn = idIn
	o.hasIDIn = true
	return o
}

func (o *searchValueQueryImplementation) SourceReferenceID() string {
	return o.sourceReferenceID
}

func (o *searchValueQueryImplementation) HasSourceReferenceID() bool {
	return o.hasSourceReferenceID
}

func (o *searchValueQueryImplementation) SetSourceReferenceID(sourceReferenceID string) SearchValueQueryInterface {
	o.sourceReferenceID = sourceReferenceID
	o.hasSourceReferenceID = true
	return o
}

func (o *searchValueQueryImplementation) SearchValue() string {
	return o.searchValue
}

func (o *searchValueQueryImplementation) HasSearchValue() bool {
	return o.hasSearchValue
}

func (o *searchValueQueryImplementation) SetSearchValue(searchValue string) SearchValueQueryInterface {
	o.searchValue = searchValue
	o.hasSearchValue = true
	return o
}

func (o *searchValueQueryImplementation) SearchType() string {
	return o.searchType
}

func (o *searchValueQueryImplementation) HasSearchType() bool {
	return o.hasSearchType
}

func (o *searchValueQueryImplementation) SetSearchType(searchType string) SearchValueQueryInterface {
	o.searchType = searchType
	o.hasSearchType = true
	return o
}

func (o *searchValueQueryImplementation) Offset() int {
	return o.offset
}

func (o *searchValueQueryImplementation) HasOffset() bool {
	return o.hasOffset
}

func (o *searchValueQueryImplementation) SetOffset(offset int) SearchValueQueryInterface {
	o.offset = offset
	o.hasOffset = true
	return o
}

func (o *searchValueQueryImplementation) Limit() int {
	return o.limit
}

func (o *searchValueQueryImplementation) HasLimit() bool {
	return o.hasLimit
}

func (o *searchValueQueryImplementation) SetLimit(limit int) SearchValueQueryInterface {
	o.limit = limit
	o.hasLimit = true
	return o
}

func (o *searchValueQueryImplementation) OrderBy() string {
	return o.orderBy
}

func (o *searchValueQueryImplementation) HasOrderBy() bool {
	return o.hasOrderBy
}

func (o *searchValueQueryImplementation) SetOrderBy(orderBy string) SearchValueQueryInterface {
	o.orderBy = orderBy
	o.hasOrderBy = true
	return o
}

func (o *searchValueQueryImplementation) OrderDirection() string {
	return o.orderDirection
}

func (o *searchValueQueryImplementation) HasOrderDirection() bool {
	return o.hasOrderDirection
}

func (o *searchValueQueryImplementation) SetOrderDirection(orderDirection string) SearchValueQueryInterface {
	o.orderDirection = orderDirection
	o.hasOrderDirection = true
	return o
}

func (o *searchValueQueryImplementation) CountOnly() bool {
	return o.countOnly
}

func (o *searchValueQueryImplementation) HasCountOnly() bool {
	return o.hasCountOnly
}

func (o *searchValueQueryImplementation) SetCountOnly(countOnly bool) SearchValueQueryInterface {
	o.countOnly = countOnly
	o.hasCountOnly = true
	return o
}

func (o *searchValueQueryImplementation) WithSoftDeleted() bool {
	return o.withSoftDeleted
}

func (o *searchValueQueryImplementation) HasWithSoftDeleted() bool {
	return o.hasWithSoftDeleted
}

func (o *searchValueQueryImplementation) SetWithSoftDeleted(withSoftDeleted bool) SearchValueQueryInterface {
	o.withSoftDeleted = withSoftDeleted
	o.hasWithSoftDeleted = true
	return o
}

func (o *searchValueQueryImplementation) Columns() []any {
	return o.columns
}

func (o *searchValueQueryImplementation) HasColumns() bool {
	return o.hasColumns
}

func (o *searchValueQueryImplementation) SetColumns(columns []any) SearchValueQueryInterface {
	o.columns = columns
	o.hasColumns = true
	return o
}

func (o *searchValueQueryImplementation) Validate() error {
	if o.hasID && strings.TrimSpace(o.id) == "" {
		return errors.New("id cannot be empty when set")
	}

	if o.hasIDIn {
		if len(o.idIn) == 0 {
			return errors.New("id_in must contain at least one value")
		}

		for _, id := range o.idIn {
			if strings.TrimSpace(id) == "" {
				return errors.New("id_in cannot contain empty values")
			}
		}
	}

	if o.hasSourceReferenceID && strings.TrimSpace(o.sourceReferenceID) == "" {
		return errors.New("source_reference_id cannot be empty when set")
	}

	if o.hasSearchValue && strings.TrimSpace(o.searchValue) == "" {
		return errors.New("search_value cannot be empty when set")
	}

	if o.hasSearchType {
		normalized := strings.ToLower(o.searchType)
		switch normalized {
		case SEARCH_TYPE_EQUALS, SEARCH_TYPE_CONTAINS, SEARCH_TYPE_STARTS_WITH, SEARCH_TYPE_ENDS_WITH:
		default:
			return errors.New("search_type must be one of equals, contains, starts_with, ends_with")
		}
	}

	if o.hasLimit && o.limit < 0 {
		return errors.New("limit cannot be negative")
	}

	if o.hasOffset && o.offset < 0 {
		return errors.New("offset cannot be negative")
	}

	if o.hasOrderBy && strings.TrimSpace(o.orderBy) == "" {
		return errors.New("order_by cannot be empty when set")
	}

	if o.hasOrderDirection {
		if !o.hasOrderBy {
			return errors.New("order_direction requires order_by to be set")
		}

		if !strings.EqualFold(o.orderDirection, "ASC") && !strings.EqualFold(o.orderDirection, "DESC") {
			return errors.New("order_direction must be asc or desc")
		}
	}

	return nil
}
