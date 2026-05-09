package blindindexstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS ==================================================================

type searchValueImplementation struct {
	dataobject.DataObject
}

var _ SearchValueInterface = (*searchValueImplementation)(nil) // verify it extends the interface

// == CONSTRUCTORS ===========================================================

func NewSearchValue() SearchValueInterface {
	d := (&searchValueImplementation{}).
		SetID(uid.HumanUid()).
		SetSourceReferenceID("").
		SetSearchValue("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	return d
}

func NewSearchValueFromExistingData(data map[string]string) SearchValueInterface {
	o := &searchValueImplementation{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================
func (d *searchValueImplementation) CreatedAt() string {
	return d.Get(COLUMN_CREATED_AT)
}

func (d *searchValueImplementation) CreatedAtCarbon() *carbon.Carbon {
	createdAt := d.CreatedAt()
	return carbon.Parse(createdAt)
}

func (d *searchValueImplementation) SetCreatedAt(createdAt string) SearchValueInterface {
	d.Set(COLUMN_CREATED_AT, createdAt)
	return d
}

func (d *searchValueImplementation) SoftDeletedAt() string {
	return d.Get(COLUMN_SOFT_DELETED_AT)
}

func (d *searchValueImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	softDeletedAt := d.SoftDeletedAt()
	return carbon.Parse(softDeletedAt)
}

func (d *searchValueImplementation) SetSoftDeletedAt(softDeletedAt string) SearchValueInterface {
	d.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return d
}

// ID returns the ID of the exam
func (o *searchValueImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the exam
func (o *searchValueImplementation) SetID(id string) SearchValueInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (d *searchValueImplementation) SourceReferenceID() string {
	return d.Get(COLUMN_SOURCE_REFERENCE_ID)
}

func (d *searchValueImplementation) SetSourceReferenceID(objectID string) SearchValueInterface {
	d.Set(COLUMN_SOURCE_REFERENCE_ID, objectID)
	return d
}

func (d *searchValueImplementation) SearchValue() string {
	return d.Get(COLUMN_SEARCH_VALUE)
}

func (d *searchValueImplementation) SetSearchValue(value string) SearchValueInterface {
	d.Set(COLUMN_SEARCH_VALUE, value)
	return d
}

func (d *searchValueImplementation) UpdatedAt() string {
	return d.Get(COLUMN_UPDATED_AT)
}

func (d *searchValueImplementation) UpdatedAtCarbon() *carbon.Carbon {
	updatedAt := d.UpdatedAt()
	return carbon.Parse(updatedAt)
}

func (d *searchValueImplementation) SetUpdatedAt(updatedAt string) SearchValueInterface {
	d.Set(COLUMN_UPDATED_AT, updatedAt)
	return d
}

func (o *searchValueImplementation) Metas() (map[string]string, error) {
	metasString := o.Get(COLUMN_METAS)

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

	o.Set(COLUMN_METAS, string(json))
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
	o.DataObject.MarkAsNotDirty()
}
