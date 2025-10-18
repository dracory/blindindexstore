package blindindexstore

import (
	"errors"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

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

	metaSearch    []string
	metaSearchNot []string

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

func (o *searchValueQueryImplementation) AddMetaSearch(needle string) SearchValueQueryInterface {
	if o.metaSearch == nil {
		o.metaSearch = []string{}
	}

	o.metaSearch = append(o.metaSearch, needle)
	return o
}

func (o *searchValueQueryImplementation) GetMetaSearch() []string {
	return o.metaSearch
}

func (o *searchValueQueryImplementation) AddMetaSearchNot(needle string) SearchValueQueryInterface {
	if o.metaSearchNot == nil {
		o.metaSearchNot = []string{}
	}

	o.metaSearchNot = append(o.metaSearchNot, needle)
	return o
}

func (o *searchValueQueryImplementation) GetMetaSearchNot() []string {
	return o.metaSearchNot
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

		if !strings.EqualFold(o.orderDirection, sb.ASC) && !strings.EqualFold(o.orderDirection, sb.DESC) {
			return errors.New("order_direction must be asc or desc")
		}
	}

	return nil
}

func (o *searchValueQueryImplementation) ToSelectDataset(store StoreInterface) (*goqu.SelectDataset, []any, error) {
	if err := o.Validate(); err != nil {
		return nil, nil, err
	}

	impl, ok := store.(*storeImplementation)
	if !ok {
		return nil, nil, errors.New("store must be a *storeImplementation")
	}

	dataset := goqu.Dialect(impl.dbDriverName).From(impl.tableName)

	if !o.HasOrderBy() || o.OrderBy() == "" {
		o.SetOrderBy(COLUMN_CREATED_AT)
	}

	if o.HasID() {
		dataset = dataset.Where(goqu.C(COLUMN_ID).Eq(o.ID()))
	}

	if o.HasIDIn() {
		dataset = dataset.Where(goqu.C(COLUMN_ID).In(o.IDIn()))
	}

	if o.HasSourceReferenceID() {
		dataset = dataset.Where(goqu.C(COLUMN_SOURCE_REFERENCE_ID).Eq(o.SourceReferenceID()))
	}

	if o.HasSearchValue() {
		searchValue := o.SearchValue()
		if impl.transformer != nil {
			searchValue = impl.transformer.Transform(searchValue)
		}

		searchType := o.SearchType()
		if !o.HasSearchType() || searchType == "" {
			searchType = SEARCH_TYPE_EQUALS
		}

		switch {
		case strings.EqualFold(searchType, SEARCH_TYPE_CONTAINS):
			dataset = dataset.Where(goqu.C(COLUMN_SEARCH_VALUE).Like("%" + searchValue + "%"))
		case strings.EqualFold(searchType, SEARCH_TYPE_STARTS_WITH):
			dataset = dataset.Where(goqu.C(COLUMN_SEARCH_VALUE).Like(searchValue + "%"))
		case strings.EqualFold(searchType, SEARCH_TYPE_ENDS_WITH):
			dataset = dataset.Where(goqu.C(COLUMN_SEARCH_VALUE).Like("%" + searchValue))
		default:
			dataset = dataset.Where(goqu.C(COLUMN_SEARCH_VALUE).Eq(searchValue))
		}
	}

	if len(o.metaSearch) > 0 {
		ors := make([]goqu.Expression, 0, len(o.metaSearch))
		for _, needle := range o.metaSearch {
			ors = append(ors, goqu.C(COLUMN_METAS).Like("%"+needle+"%"))
		}

		if len(ors) > 0 {
			dataset = dataset.Where(goqu.Or(ors...))
		}
	}

	if len(o.metaSearchNot) > 0 {
		for _, needle := range o.metaSearchNot {
			dataset = dataset.Where(goqu.C(COLUMN_METAS).NotLike("%" + needle + "%"))
		}
	}

	countOnly := o.CountOnly()
	if !o.HasCountOnly() {
		countOnly = false
	}

	if !countOnly {
		if o.HasLimit() && o.Limit() > 0 {
			dataset = dataset.Limit(uint(o.Limit()))
		}

		if o.HasOffset() && o.Offset() > 0 {
			dataset = dataset.Offset(uint(o.Offset()))
		}
	}

	sortOrder := o.OrderDirection()
	if !o.HasOrderDirection() || sortOrder == "" {
		sortOrder = sb.DESC
	}

	if o.HasOrderBy() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			dataset = dataset.Order(goqu.I(o.OrderBy()).Asc())
		} else {
			dataset = dataset.Order(goqu.I(o.OrderBy()).Desc())
		}
	}

	includeSoftDeleted := o.HasWithSoftDeleted() && o.WithSoftDeleted()
	if !includeSoftDeleted {
		dataset = dataset.Where(goqu.C(COLUMN_SOFT_DELETED_AT).Gt(carbon.Now(carbon.UTC).ToDateTimeString()))
	}

	columns := o.Columns()
	if !o.HasColumns() || len(columns) == 0 {
		columns = []any{"*"}
	}

	return dataset, columns, nil
}
