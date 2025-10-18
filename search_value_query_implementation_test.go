package blindindexstore

import (
	"strings"
	"testing"
)

func TestSearchValueQueryToSelectDataset_Defaults(t *testing.T) {
	store := &storeImplementation{
		tableName:    "test_search_values",
		dbDriverName: "sqlite",
	}

	query := NewSearchValueQuery().SetSearchValue("needle")

	dataset, columns, err := query.ToSelectDataset(store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 1 {
		t.Fatalf("expected exactly one column, got %d", len(columns))
	}

	if columns[0] != "*" {
		t.Fatalf("expected default column '*', got %v", columns[0])
	}

	sqlStr, params, err := dataset.Select(columns...).Prepared(true).ToSQL()
	if err != nil {
		t.Fatalf("unexpected error generating SQL: %v", err)
	}

	if !strings.Contains(sqlStr, "ORDER BY \"created_at\" DESC") {
		t.Fatalf("expected SQL to order by created_at DESC, got %s", sqlStr)
	}

	if !strings.Contains(sqlStr, "\"soft_deleted_at\" > ?") {
		t.Fatalf("expected SQL to filter by soft_deleted_at, got %s", sqlStr)
	}

	if len(params) != 2 {
		t.Fatalf("expected two SQL parameters, got %d", len(params))
	}

	if params[0] != "needle" {
		t.Fatalf("expected first parameter to equal search value, got %v", params[0])
	}
}
