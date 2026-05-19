package blindindexstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/doug-martin/goqu/v9"

	"github.com/dracory/database"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// storeImplementation implements StoreInterface
type storeImplementation struct {
	tableName          string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
	transformer        TransformerInterface
}

// AutoMigrate auto migrate (deprecated - use MigrateUp)
func (st *storeImplementation) AutoMigrate() error {
	return st.MigrateUp()
}

// MigrateUp creates the table
func (st *storeImplementation) MigrateUp(tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	sqlStr, err := st.sqlTableCreate()
	if err != nil {
		return err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	var errExec error
	if txToUse != nil {
		_, errExec = txToUse.Exec(sqlStr)
	} else {
		_, errExec = st.db.Exec(sqlStr)
	}

	if errExec != nil {
		log.Println(errExec)
		return errExec
	}

	return nil
}

// MigrateDown drops the table
func (st *storeImplementation) MigrateDown(tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	sqlStr, err := st.sqlTableDrop()
	if err != nil {
		return err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	var errExec error
	if txToUse != nil {
		_, errExec = txToUse.Exec(sqlStr)
	} else {
		_, errExec = st.db.Exec(sqlStr)
	}

	if errExec != nil {
		log.Println(errExec)
		return errExec
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
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

	// Build the select dataset
	ds, _, err := query.ToSelectDataset(store)
	if err != nil {
		return nil, err
	}

	// Generate SQL and parameters
	sqlStr, sqlParams, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(database.NewQueryableContext(ctx, store.db), sqlStr, sqlParams...)
	if err != nil {
		return refIDs, err
	}

	list := []string{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewSearchValueFromExistingData(modelMap)
		list = append(list, model.SourceReferenceID())
	})

	return list, nil
}

// SearchValueCreate creates the record
// Side effect! Transforms the value
func (store *storeImplementation) SearchValueCreate(ctx context.Context, searchValue SearchValueInterface) error {
	searchValue.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	searchValue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	searchValue.SetSearchValue(store.transformer.Transform(searchValue.SearchValue()))

	data := searchValue.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.tableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.tableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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
	// Build the select dataset
	ds, _, err := query.ToSelectDataset(store)
	if err != nil {
		return nil, err
	}

	// Generate SQL and parameters
	sqlStr, sqlParams, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(database.NewQueryableContext(ctx, store.db), sqlStr, sqlParams...)
	if err != nil {
		return []SearchValueInterface{}, err
	}

	list := []SearchValueInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewSearchValueFromExistingData(modelMap)
		list = append(list, model)
	})

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

	dataChanged := searchValue.DataChanged()

	delete(dataChanged, "id")   // ID is not updateable
	delete(dataChanged, "hash") // Hash is not updateable
	delete(dataChanged, "data") // Data is not updateable

	if len(dataChanged) < 2 {
		return nil
	}

	if lo.HasKey(dataChanged, COLUMN_SEARCH_VALUE) {
		searchValue.SetSearchValue(store.transformer.Transform(searchValue.SearchValue()))
		dataChanged[COLUMN_SEARCH_VALUE] = searchValue.SearchValue()
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.tableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(searchValue.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	searchValue.MarkAsNotDirty()

	return err
}

// IsAutomigrateEnabled returns whether automigrate is enabled
func (st *storeImplementation) IsAutomigrateEnabled() bool {
	return st.automigrateEnabled
}

func (store *storeImplementation) Truncate(ctx context.Context) error {
	var (
		sqlStr string
		errSql error
	)

	if strings.EqualFold(store.dbDriverName, "sqlite") {
		sqlStr, _, errSql = goqu.Dialect(store.dbDriverName).
			Delete(store.tableName).
			Prepared(true).
			ToSQL()
	} else {
		sqlStr, _, errSql = goqu.Dialect(store.dbDriverName).
			Truncate(store.tableName).
			ToSQL()
	}

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr)

	return err
}

// 	}

// 	if options.SourceReferenceID != "" {
// 		q = q.Where(goqu.C(COLUMN_SOURCE_REFERENCE_ID).Eq(options.SourceReferenceID))
// 	}

// 	if options.SearchValue != "" {
// 		options.SearchValue = store.transformer.Transform(options.SearchValue)
// 		if options.SearchType == SEARCH_TYPE_CONTAINS {
// 			q = q.Where(goqu.C(COLUMN_SEARCH_VALUE).Like("%" + options.SearchValue + "%"))
// 		} else if options.SearchType == SEARCH_TYPE_STARTS_WITH {
// 			q = q.Where(goqu.C(COLUMN_SEARCH_VALUE).Like(options.SearchValue + "%"))
// 		} else if options.SearchType == SEARCH_TYPE_ENDS_WITH {
// 			q = q.Where(goqu.C(COLUMN_SEARCH_VALUE).Like(options.SearchValue + "%"))
// 		} else {
// 			// default to strict search
// 			q = q.Where(goqu.C(COLUMN_SEARCH_VALUE).Eq(options.SearchValue))
// 		}
// 	}

// 	if !options.CountOnly {
// 		if options.Limit > 0 {
// 			q = q.Limit(uint(options.Limit))
// 		}

// 		if options.Offset > 0 {
// 			q = q.Offset(uint(options.Offset))
// 		}
// 	}

// 	sortOrder := sb.DESC
// 	if options.SortOrder != "" {
// 		sortOrder = options.SortOrder
// 	}

// 	if options.OrderBy != "" {
// 		if strings.EqualFold(sortOrder, sb.ASC) {
// 			q = q.Order(goqu.I(options.OrderBy).Asc())
// 		} else {
// 			q = q.Order(goqu.I(options.OrderBy).Desc())
// 		}
// 	}

// 	if !options.WithSoftDeleted {
// 		q = q.Where(goqu.C(COLUMN_SOFT_DELETED_AT).Gt(carbon.Now(carbon.UTC).ToDateTimeString()))
// 	}

// 	return q
// }
