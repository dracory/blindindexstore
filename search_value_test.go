package blindindexstore

import (
	"testing"

	"github.com/dracory/sb"
)

func TestNewSearchValueDefaults(t *testing.T) {
	value := NewSearchValue()

	if value == nil {
		t.Fatal("NewSearchValue returned nil")
	}

	if value.ID() == "" {
		t.Fatal("ID should not be empty")
	}

	if value.SourceReferenceID() != "" {
		t.Fatalf("expected empty source reference ID, got %q", value.SourceReferenceID())
	}

	if value.SearchValue() != "" {
		t.Fatalf("expected empty search value, got %q", value.SearchValue())
	}

	if value.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatalf("expected soft deleted at to equal sb.MAX_DATETIME, got %q", value.SoftDeletedAt())
	}

	if value.CreatedAt() == "" {
		t.Fatal("CreatedAt should not be empty")
	}

	if value.UpdatedAt() == "" {
		t.Fatal("UpdatedAt should not be empty")
	}
}

func TestNewSearchValueFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:                 "test-id",
		COLUMN_SOURCE_REFERENCE_ID: "source-1",
		COLUMN_SEARCH_VALUE:       "plain-value",
		COLUMN_SOFT_DELETED_AT:    "2099-01-01 00:00:00",
		COLUMN_CREATED_AT:         "2023-01-01 01:02:03",
		COLUMN_UPDATED_AT:         "2023-01-02 04:05:06",
		COLUMN_METAS:              `{"foo":"bar"}`,
	}

	value := NewSearchValueFromExistingData(data)

	if value.ID() != data[COLUMN_ID] {
		t.Fatalf("expected ID %q, got %q", data[COLUMN_ID], value.ID())
	}

	if value.SourceReferenceID() != data[COLUMN_SOURCE_REFERENCE_ID] {
		t.Fatalf("expected source reference ID %q, got %q", data[COLUMN_SOURCE_REFERENCE_ID], value.SourceReferenceID())
	}

	if value.SearchValue() != data[COLUMN_SEARCH_VALUE] {
		t.Fatalf("expected search value %q, got %q", data[COLUMN_SEARCH_VALUE], value.SearchValue())
	}

	if value.SoftDeletedAt() != data[COLUMN_SOFT_DELETED_AT] {
		t.Fatalf("expected soft deleted at %q, got %q", data[COLUMN_SOFT_DELETED_AT], value.SoftDeletedAt())
	}

	if value.CreatedAt() != data[COLUMN_CREATED_AT] {
		t.Fatalf("expected created at %q, got %q", data[COLUMN_CREATED_AT], value.CreatedAt())
	}

	if value.UpdatedAt() != data[COLUMN_UPDATED_AT] {
		t.Fatalf("expected updated at %q, got %q", data[COLUMN_UPDATED_AT], value.UpdatedAt())
	}

	metas, err := value.Metas()
	if err != nil {
		t.Fatalf("unexpected error reading metas: %v", err)
	}

	if len(metas) != 1 {
		t.Fatalf("expected 1 meta entry, got %d", len(metas))
	}

	if metas["foo"] != "bar" {
		t.Fatalf("expected meta foo to be bar, got %q", metas["foo"])
	}
}

func TestSearchValueMetasLifecycle(t *testing.T) {
	value := NewSearchValue()

	metas, err := value.Metas()
	if err != nil {
		t.Fatalf("unexpected error reading metas: %v", err)
	}

	if metas != nil {
		t.Fatalf("expected metas to be nil, got %v", metas)
	}

	has, err := value.HasMeta("foo")
	if err != nil {
		t.Fatalf("unexpected error checking meta: %v", err)
	}

	if has {
		t.Fatal("expected has meta to be false before setting")
	}

	value, err = value.SetMeta("foo", "bar")
	if err != nil {
		t.Fatalf("unexpected error setting meta: %v", err)
	}

	has, err = value.HasMeta("foo")
	if err != nil {
		t.Fatalf("unexpected error checking meta after set: %v", err)
	}

	if !has {
		t.Fatal("expected has meta to be true after setting")
	}

	val, err := value.Meta("foo")
	if err != nil {
		t.Fatalf("unexpected error retrieving meta: %v", err)
	}

	if val != "bar" {
		t.Fatalf("expected meta foo to equal bar, got %q", val)
	}

	value, err = value.SetMetas(map[string]string{"alpha": "1", "beta": "2"})
	if err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err = value.Metas()
	if err != nil {
		t.Fatalf("unexpected error reading metas after set: %v", err)
	}

	if len(metas) != 2 {
		t.Fatalf("expected 2 meta entries, got %d", len(metas))
	}

	value, err = value.DeleteMeta("alpha")
	if err != nil {
		t.Fatalf("unexpected error deleting meta: %v", err)
	}

	has, err = value.HasMeta("alpha")
	if err != nil {
		t.Fatalf("unexpected error checking meta after delete: %v", err)
	}

	if has {
		t.Fatal("expected has meta to be false after delete")
	}

	val, err = value.Meta("alpha")
	if err != nil {
		t.Fatalf("unexpected error retrieving meta after delete: %v", err)
	}

	if val != "" {
		t.Fatalf("expected empty meta after delete, got %q", val)
	}
}

func TestSearchValueSetters(t *testing.T) {
	value := NewSearchValue()

	expectedID := "id-123"
	value = value.SetID(expectedID)
	if value.ID() != expectedID {
		t.Fatalf("expected ID %q, got %q", expectedID, value.ID())
	}

	expectedRef := "ref-456"
	value = value.SetSourceReferenceID(expectedRef)
	if value.SourceReferenceID() != expectedRef {
		t.Fatalf("expected source reference ID %q, got %q", expectedRef, value.SourceReferenceID())
	}

	expectedSearch := "search-value"
	value = value.SetSearchValue(expectedSearch)
	if value.SearchValue() != expectedSearch {
		t.Fatalf("expected search value %q, got %q", expectedSearch, value.SearchValue())
	}

	expectedCreated := "2024-01-02 03:04:05"
	value = value.SetCreatedAt(expectedCreated)
	if value.CreatedAt() != expectedCreated {
		t.Fatalf("expected created at %q, got %q", expectedCreated, value.CreatedAt())
	}

	expectedUpdated := "2024-01-05 06:07:08"
	value = value.SetUpdatedAt(expectedUpdated)
	if value.UpdatedAt() != expectedUpdated {
		t.Fatalf("expected updated at %q, got %q", expectedUpdated, value.UpdatedAt())
	}

	expectedSoft := "2025-02-03 04:05:06"
	value = value.SetSoftDeletedAt(expectedSoft)
	if value.SoftDeletedAt() != expectedSoft {
		t.Fatalf("expected soft deleted at %q, got %q", expectedSoft, value.SoftDeletedAt())
	}
}
