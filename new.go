package blindindexstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/neat"
)

// == OPTIONS ==================================================================

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	TableName          string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
	Transformer        TransformerInterface
}

// == CONSTRUCTOR =============================================================

// NewStore creates a new entity store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.DB == nil {
		return nil, errors.New("blind index store: DB is required")
	}

	if opts.TableName == "" {
		return nil, errors.New("blind index store: TableName is required")
	}

	if opts.Transformer == nil {
		return nil, errors.New("blind index store: Transformer is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	store := &storeImplementation{
		tableName:          opts.TableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 neatDB,
		debugEnabled:       opts.DebugEnabled,
		transformer:        opts.Transformer,
	}

	if store.automigrateEnabled {
		err := store.MigrateUp(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
