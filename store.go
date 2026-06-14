package blindindexstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
)

// == INTERFACE ===============================================================

type StoreInterface interface {
	// GetTableName returns the table name
	GetTableName() string
	// SetTableName sets the table name
	SetTableName(tableName string)

	// MigrateDown drops the table
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	// MigrateUp creates the table
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	EnableDebug(debug bool)
	Search(ctx context.Context, needle, searchType string) (refIDs []string, err error)
	SearchValueCreate(ctx context.Context, value SearchValueInterface) error
	SearchValueDelete(ctx context.Context, value SearchValueInterface) error
	SearchValueDeleteByID(ctx context.Context, valueID string) error
	SearchValueFindByID(ctx context.Context, id string) (SearchValueInterface, error)
	SearchValueFindBySourceReferenceID(ctx context.Context, sourceReferenceID string) (SearchValueInterface, error)
	SearchValueList(ctx context.Context, query SearchValueQueryInterface) ([]SearchValueInterface, error)
	SearchValueSoftDelete(ctx context.Context, discount SearchValueInterface) error
	SearchValueSoftDeleteByID(ctx context.Context, discountID string) error
	SearchValueUpdate(ctx context.Context, value SearchValueInterface) error
	Truncate(ctx context.Context) error

	// IsAutomigrateEnabled returns whether automigrate is enabled
	IsAutomigrateEnabled() bool
}

// == TYPE ====================================================================

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// storeImplementation implements StoreInterface
type storeImplementation struct {
	tableName          string
	db                 *neat.Database
	automigrateEnabled bool
	debugEnabled       bool
	transformer        TransformerInterface
}

// AutoMigrate auto migrate (deprecated - use MigrateUp)
func (st *storeImplementation) AutoMigrate() error {
	return st.MigrateUp(context.Background())
}

// MigrateUp creates the table
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.tableName) {
		if st.debugEnabled {
			log.Println("MigrateUp: table already exists", "table", st.tableName)
		}
		return nil
	}

	err := st.db.Schema().Create(st.tableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_SOURCE_REFERENCE_ID, 40)
		table.Text(COLUMN_SEARCH_VALUE)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})

	if err != nil {
		if st.debugEnabled {
			log.Println("MigrateUp failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the table
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if !st.db.Schema().HasTable(st.tableName) {
		if st.debugEnabled {
			log.Println("MigrateDown: table does not exist", "table", st.tableName)
		}
		return nil
	}

	err := st.db.Schema().Drop(st.tableName)
	if err != nil {
		if st.debugEnabled {
			log.Println("MigrateDown failed", "error", err)
		}
		return err
	}
	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if st.db != nil {
		if debug {
			st.db.EnableDebug()
		} else {
			st.db.DisableDebug()
		}
	}
}

// GetTableName returns the table name
func (st *storeImplementation) GetTableName() string {
	return st.tableName
}

// SetTableName sets the table name
func (st *storeImplementation) SetTableName(tableName string) {
	st.tableName = tableName
}

func (store *storeImplementation) Search(ctx context.Context, needle, searchType string) (refIDs []string, err error) {
	query := NewSearchValueQuery().
		SetSearchValue(needle).
		SetSearchType(searchType)

	q := store.buildQuery(query)

	type searchValueRow struct {
		ID                string `db:"id"`
		SourceReferenceID string `db:"source_reference_id"`
	}

	var rows []searchValueRow
	if err := q.Table(store.tableName).Get(&rows); err != nil {
		return []string{}, err
	}

	list := []string{}
	for _, r := range rows {
		list = append(list, r.SourceReferenceID)
	}

	return list, nil
}

// SearchValueCreate creates the record
// Side effect! Transforms the value
func (store *storeImplementation) SearchValueCreate(ctx context.Context, searchValue SearchValueInterface) error {
	searchValue.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	searchValue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	searchValue.SetSearchValue(store.transformer.Transform(searchValue.SearchValue()))

	row := map[string]any{
		COLUMN_ID:                  searchValue.ID(),
		COLUMN_SOURCE_REFERENCE_ID: searchValue.SourceReferenceID(),
		COLUMN_SEARCH_VALUE:        searchValue.SearchValue(),
		COLUMN_CREATED_AT:          searchValue.CreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:          searchValue.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:     searchValue.SoftDeletedAtCarbon().StdTime(),
	}

	err := store.db.Query().Table(store.tableName).Create(row)
	if err != nil {
		return err
	}

	searchValue.MarkAsNotDirty()

	return nil
}

func (store *storeImplementation) SearchValueDelete(ctx context.Context, searchValue SearchValueInterface) error {
	if searchValue == nil {
		return errors.New("searchValue is nil")
	}

	return store.SearchValueDeleteByID(ctx, searchValue.ID())
}

func (store *storeImplementation) SearchValueDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("searchValue id is empty")
	}

	_, err := store.db.Query().
		Table(store.tableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

func (store *storeImplementation) SearchValueFindByID(ctx context.Context, id string) (SearchValueInterface, error) {
	if id == "" {
		return nil, errors.New("searchValue id is empty")
	}

	list, err := store.SearchValueList(ctx, NewSearchValueQuery().
		SetID(id).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) SearchValueFindBySourceReferenceID(ctx context.Context, sourceReferenceID string) (SearchValueInterface, error) {
	if sourceReferenceID == "" {
		return nil, errors.New("searchValue objectID is empty")
	}

	list, err := store.SearchValueList(ctx, NewSearchValueQuery().
		SetSourceReferenceID(sourceReferenceID).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) SearchValueList(ctx context.Context, query SearchValueQueryInterface) ([]SearchValueInterface, error) {
	q := store.buildQuery(query)

	type searchValueRow struct {
		ID                string     `db:"id"`
		SourceReferenceID string     `db:"source_reference_id"`
		SearchValue       string     `db:"search_value"`
		CreatedAt         time.Time  `db:"created_at"`
		UpdatedAt         time.Time  `db:"updated_at"`
		SoftDeletedAt     *time.Time `db:"soft_deleted_at"`
	}

	var rows []searchValueRow
	if err := q.Table(store.tableName).Get(&rows); err != nil {
		return []SearchValueInterface{}, err
	}

	list := []SearchValueInterface{}
	for _, r := range rows {
		s := NewSearchValue()
		s.SetID(r.ID)
		s.SetSourceReferenceID(r.SourceReferenceID)
		s.SetSearchValue(r.SearchValue)
		s.SetCreatedAt(carbon.CreateFromStdTime(r.CreatedAt).ToDateTimeString())
		s.SetUpdatedAt(carbon.CreateFromStdTime(r.UpdatedAt).ToDateTimeString())
		if r.SoftDeletedAt != nil {
			s.SetSoftDeletedAt(carbon.CreateFromStdTime(*r.SoftDeletedAt).ToDateTimeString())
		} else {
			s.SetSoftDeletedAt(MAX_DATETIME)
		}
		list = append(list, s)
	}

	return list, nil
}

func (store *storeImplementation) SearchValueSoftDelete(ctx context.Context, searchValue SearchValueInterface) error {
	if searchValue == nil {
		return errors.New("searchValue is nil")
	}

	searchValue.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SearchValueUpdate(ctx, searchValue)
}

func (store *storeImplementation) SearchValueSoftDeleteByID(ctx context.Context, id string) error {
	searchValue, err := store.SearchValueFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.SearchValueSoftDelete(ctx, searchValue)
}

// SearchValueUpdate updates the record
// Side effect! Transforms the value, use with caution
func (store *storeImplementation) SearchValueUpdate(ctx context.Context, searchValue SearchValueInterface) error {
	if searchValue == nil {
		return errors.New("order is nil")
	}

	searchValue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	row := map[string]any{
		COLUMN_UPDATED_AT: searchValue.UpdatedAtCarbon().StdTime(),
	}

	if searchValue.SourceReferenceID() != "" {
		row[COLUMN_SOURCE_REFERENCE_ID] = searchValue.SourceReferenceID()
	}

	if searchValue.SearchValue() != "" {
		// Only transform if the value appears to be plain text (not already hashed)
		// SHA256 produces 64-character hex strings
		currentValue := searchValue.SearchValue()
		var transformedValue string
		if len(currentValue) == 64 && isHex(currentValue) {
			// Value appears to be already hashed, use as-is
			transformedValue = currentValue
		} else {
			// Value appears to be plain text, transform it
			transformedValue = store.transformer.Transform(currentValue)
		}
		row[COLUMN_SEARCH_VALUE] = transformedValue
		searchValue.SetSearchValue(transformedValue)
	}

	if searchValue.SoftDeletedAt() != "" {
		row[COLUMN_SOFT_DELETED_AT] = searchValue.SoftDeletedAtCarbon().StdTime()
	}

	_, err := store.db.Query().
		Table(store.tableName).
		Where(COLUMN_ID+" = ?", searchValue.ID()).
		Update(row)

	searchValue.MarkAsNotDirty()

	return err
}

// IsAutomigrateEnabled returns whether automigrate is enabled
func (st *storeImplementation) IsAutomigrateEnabled() bool {
	return st.automigrateEnabled
}

func (store *storeImplementation) Truncate(ctx context.Context) error {
	_, err := store.db.Query().
		Table(store.tableName).
		Delete()
	return err
}

// buildQuery builds a neat query from the search value query interface.
func (st *storeImplementation) buildQuery(query SearchValueQueryInterface) contractsorm.Query {
	q := st.db.Query()

	if query == nil {
		return q
	}

	if query.HasID() && query.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.ID())
	}

	if query.HasIDIn() && len(query.IDIn()) > 0 {
		q = q.Where(COLUMN_ID+" IN (?)", query.IDIn())
	}

	if query.HasSourceReferenceID() && query.SourceReferenceID() != "" {
		q = q.Where(COLUMN_SOURCE_REFERENCE_ID+" = ?", query.SourceReferenceID())
	}

	if query.HasSearchValue() && query.SearchValue() != "" {
		searchValue := st.transformer.Transform(query.SearchValue())
		searchType := query.SearchType()
		if !query.HasSearchType() || searchType == "" {
			searchType = SEARCH_TYPE_EQUALS
		}

		switch searchType {
		case SEARCH_TYPE_CONTAINS:
			q = q.Where(COLUMN_SEARCH_VALUE+" LIKE ?", "%"+searchValue+"%")
		case SEARCH_TYPE_STARTS_WITH:
			q = q.Where(COLUMN_SEARCH_VALUE+" LIKE ?", searchValue+"%")
		case SEARCH_TYPE_ENDS_WITH:
			q = q.Where(COLUMN_SEARCH_VALUE+" LIKE ?", "%"+searchValue)
		default:
			q = q.Where(COLUMN_SEARCH_VALUE+" = ?", searchValue)
		}
	}

	if query.HasLimit() && query.Limit() > 0 {
		q = q.Limit(query.Limit())
	}

	if query.HasOffset() && query.Offset() > 0 {
		q = q.Offset(query.Offset())
	}

	if query.HasOrderBy() && query.OrderBy() != "" {
		orderDirection := query.OrderDirection()
		if !query.HasOrderDirection() || orderDirection == "" {
			orderDirection = "DESC"
		}
		if strings.EqualFold(orderDirection, "ASC") {
			q = q.OrderBy(query.OrderBy() + " ASC")
		} else {
			q = q.OrderBy(query.OrderBy() + " DESC")
		}
	} else {
		q = q.OrderBy(COLUMN_CREATED_AT + " DESC")
	}

	// Handle soft delete filtering
	if query.HasWithSoftDeleted() && query.WithSoftDeleted() {
		q = q.WithSoftDeleted()
	} else {
		// By default, filter out soft-deleted records
		q = q.Where(COLUMN_SOFT_DELETED_AT+" > ?", carbon.Now(carbon.UTC).StdTime())
	}

	return q
}
