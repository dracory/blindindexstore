package blindindexstore

import "context"

type StoreInterface interface {
	AutoMigrate() error

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
