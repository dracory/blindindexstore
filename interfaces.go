package blindindexstore

type StoreInterface interface {
	AutoMigrate() error

	Search(needle, searchType string) (refIDs []string, err error)
	SearchValueCreate(value SearchValueInterface) error
	SearchValueDelete(value SearchValueInterface) error
	SearchValueDeleteByID(valueID string) error
	SearchValueFindByID(id string) (SearchValueInterface, error)
	SearchValueFindBySourceReferenceID(sourceReferenceID string) (SearchValueInterface, error)
	SearchValueList(options SearchValueQueryOptions) ([]SearchValueInterface, error)
	SearchValueSoftDelete(discount SearchValueInterface) error
	SearchValueSoftDeleteByID(discountID string) error
	SearchValueUpdate(value SearchValueInterface) error
	Truncate() error

	// IsAutomigrateEnabled returns whether automigrate is enabled
	IsAutomigrateEnabled() bool
}

type SearchValueQueryOptions struct {
	ID                string
	IDIn              string
	SourceReferenceID string
	SearchValue       string
	SearchType        string
	Offset            int
	Limit             int
	SortOrder         string
	OrderBy           string
	CountOnly         bool
	WithSoftDeleted   bool
}
