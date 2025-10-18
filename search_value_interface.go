package blindindexstore

import "github.com/dromara/carbon/v2"

// SearchValueInterface defines the methods for a SearchValue entity
// This interface can be implemented by any SearchValue struct for flexibility and testability.
type SearchValueInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	SetCreatedAt(createdAt string) SearchValueInterface
	CreatedAtCarbon() *carbon.Carbon

	SearchValue() string
	SetSearchValue(value string) SearchValueInterface

	SourceReferenceID() string
	SetSourceReferenceID(sourceReferenceID string) SearchValueInterface

	ID() string
	SetID(id string) SearchValueInterface

	Metas() (map[string]string, error)
	SetMetas(data map[string]string) (SearchValueInterface, error)

	HasMeta(key string) (bool, error)
	Meta(key string) (string, error)
	SetMeta(key string, value string) (SearchValueInterface, error)
	DeleteMeta(key string) (SearchValueInterface, error)

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) SearchValueInterface
	SoftDeletedAtCarbon() *carbon.Carbon

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) SearchValueInterface
	UpdatedAtCarbon() *carbon.Carbon
}
