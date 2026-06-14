package blindindexstore

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dracory/neat/database/orm"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == INTERFACE ===============================================================

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

// == CLASS ==================================================================

type searchValueImplementation struct {
	orm.ShortID

	SourceReferenceIDField string    `db:"source_reference_id"`
	SearchValueField       string    `db:"search_value"`
	MetasField             string    `db:"metas"`
	CreatedAtField         time.Time `db:"created_at"`
	UpdatedAtField         time.Time `db:"updated_at"`
	orm.SoftDeletes
	DeletedAt sql.NullTime `db:"soft_deleted_at"`
}

var _ SearchValueInterface = (*searchValueImplementation)(nil) // verify it extends the interface

// == CONSTRUCTORS ===========================================================

func NewSearchValue() SearchValueInterface {
	o := &searchValueImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetSourceReferenceID("")
	o.SetSearchValue("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	return o
}

func NewSearchValueFromExistingData(data map[string]string) SearchValueInterface {
	o := &searchValueImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetSourceReferenceID(data[COLUMN_SOURCE_REFERENCE_ID])
	o.SetSearchValue(data[COLUMN_SEARCH_VALUE])
	if v, ok := data[COLUMN_METAS]; ok {
		o.MetasField = v
	}
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================
func (d *searchValueImplementation) CreatedAt() string {
	if d.CreatedAtField.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(d.CreatedAtField).ToDateTimeString()
}

func (d *searchValueImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(d.CreatedAtField)
}

func (d *searchValueImplementation) SetCreatedAt(createdAt string) SearchValueInterface {
	if createdAt == "" {
		return d
	}
	d.CreatedAtField = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return d
}

func (d *searchValueImplementation) SoftDeletedAt() string {
	if !d.DeletedAt.Valid || d.DeletedAt.Time.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(d.DeletedAt.Time).ToDateTimeString()
}

func (d *searchValueImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	if !d.DeletedAt.Valid {
		return carbon.CreateFromStdTime(time.Time{})
	}
	return carbon.CreateFromStdTime(d.DeletedAt.Time)
}

func (d *searchValueImplementation) SetSoftDeletedAt(softDeletedAt string) SearchValueInterface {
	if softDeletedAt == "" {
		d.DeletedAt = sql.NullTime{Valid: false}
		return d
	}
	t := carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	d.DeletedAt = sql.NullTime{Time: t, Valid: true}
	return d
}

// ID returns the ID of the exam
func (o *searchValueImplementation) ID() string {
	return o.ShortID.ID
}

// SetID sets the ID of the exam
func (o *searchValueImplementation) SetID(id string) SearchValueInterface {
	o.ShortID.ID = id
	return o
}

func (d *searchValueImplementation) SourceReferenceID() string {
	return d.SourceReferenceIDField
}

func (d *searchValueImplementation) SetSourceReferenceID(objectID string) SearchValueInterface {
	d.SourceReferenceIDField = objectID
	return d
}

func (d *searchValueImplementation) SearchValue() string {
	return d.SearchValueField
}

func (d *searchValueImplementation) SetSearchValue(value string) SearchValueInterface {
	d.SearchValueField = value
	return d
}

func (d *searchValueImplementation) UpdatedAt() string {
	if d.UpdatedAtField.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(d.UpdatedAtField).ToDateTimeString()
}

func (d *searchValueImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(d.UpdatedAtField)
}

func (d *searchValueImplementation) SetUpdatedAt(updatedAt string) SearchValueInterface {
	if updatedAt == "" {
		return d
	}
	d.UpdatedAtField = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return d
}

func (o *searchValueImplementation) Metas() (map[string]string, error) {
	metasString := o.MetasField

	if metasString == "" {
		return nil, nil
	}

	var metas map[string]string
	err := json.Unmarshal([]byte(metasString), &metas)

	if err != nil {
		return nil, err
	}

	return metas, nil
}

func (o *searchValueImplementation) SetMetas(data map[string]string) (SearchValueInterface, error) {
	json, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	o.MetasField = string(json)
	return o, nil
}

func (o *searchValueImplementation) HasMeta(key string) (bool, error) {
	metas, err := o.Metas()
	if err != nil {
		return false, err
	}
	if metas == nil {
		return false, nil
	}
	_, exists := metas[key]
	return exists, nil
}

func (o *searchValueImplementation) Meta(key string) (string, error) {
	metas, err := o.Metas()
	if err != nil {
		return "", err
	}
	if metas == nil {
		return "", nil
	}
	val, exists := metas[key]
	if !exists {
		return "", nil
	}
	return val, nil
}

func (o *searchValueImplementation) SetMeta(key string, value string) (SearchValueInterface, error) {
	metas, err := o.Metas()
	if err != nil {
		return nil, err
	}
	if metas == nil {
		metas = map[string]string{}
	}
	metas[key] = value
	return o.SetMetas(metas)
}

func (o *searchValueImplementation) DeleteMeta(key string) (SearchValueInterface, error) {
	metas, err := o.Metas()
	if err != nil {
		return nil, err
	}
	if metas == nil {
		return o, nil
	}
	delete(metas, key)
	return o.SetMetas(metas)
}

func (o *searchValueImplementation) MarkAsNotDirty() {
	// No-op with neat ORM traits
}

func (o *searchValueImplementation) Data() map[string]string {
	return map[string]string{
		COLUMN_ID:                  o.ID(),
		COLUMN_SOURCE_REFERENCE_ID: o.SourceReferenceID(),
		COLUMN_SEARCH_VALUE:        o.SearchValue(),
		COLUMN_METAS:               o.MetasField,
		COLUMN_CREATED_AT:          o.CreatedAt(),
		COLUMN_UPDATED_AT:          o.UpdatedAt(),
		COLUMN_SOFT_DELETED_AT:     o.SoftDeletedAt(),
	}
}

func (o *searchValueImplementation) DataChanged() map[string]string {
	// Return all fields as changed since neat ORM traits don't track dirty state
	return o.Data()
}
