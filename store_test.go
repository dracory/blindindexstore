package blindindexstore

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		panic(err)
	}

	return db
}

func Test_Store_SearchValueFindBySourceReferenceID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_find_by_source_reference_id",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	if err := store.SearchValueCreate(context.Background(), value); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.SearchValueFindBySourceReferenceID(context.Background(), "RefId01")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if found == nil || found.SourceReferenceID() != "RefId01" {
		t.Fatalf("expected to find search value by source reference id")
	}

	missing, err := store.SearchValueFindBySourceReferenceID(context.Background(), "missing")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if missing != nil {
		t.Fatalf("expected missing search value to be nil")
	}

	_, err = store.SearchValueFindBySourceReferenceID(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error for empty source reference id")
	}
}

func Test_Store_EnableDebug(t *testing.T) {
	store := &storeImplementation{}
	store.EnableDebug(true)
	if !store.debugEnabled {
		t.Fatalf("expected debugEnabled to be true")
	}
	store.EnableDebug(false)
	if store.debugEnabled {
		t.Fatalf("expected debugEnabled to be false")
	}
}

func Test_Store_WithAutoMigrate(t *testing.T) {
	db := initDB()

	storeAutomigrateFalse, errAutomigrateFalse := NewStore(NewStoreOptions{
		TableName:          "test_blindindex_with_automigrate_false",
		DB:                 db,
		AutomigrateEnabled: false,
		Transformer:        &NoChangeTransformer{},
	})

	if errAutomigrateFalse != nil {
		t.Fatalf("automigrateEnabled: Expected [err] to be nill received [%v]", errAutomigrateFalse.Error())
	}

	if storeAutomigrateFalse.IsAutomigrateEnabled() != false {
		t.Fatalf("automigrateEnabled: Expected [false] received [%v]", storeAutomigrateFalse.IsAutomigrateEnabled())
	}

	storeAutomigrateTrue, errAutomigrateTrue := NewStore(NewStoreOptions{
		TableName:          "test_blockindex_with_automigrate_true",
		DB:                 db,
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if errAutomigrateTrue != nil {
		t.Fatalf("automigrateEnabled: Expected [err] to be nill received [%v]", errAutomigrateTrue.Error())
	}

	if storeAutomigrateTrue.IsAutomigrateEnabled() != true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", storeAutomigrateTrue.IsAutomigrateEnabled())
	}

	if err := storeAutomigrateTrue.MigrateUp(context.Background()); err != nil {
		t.Fatalf("unexpected automigrate error: %v", err)
	}
}

func Test_Store_SearchValueCreate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_create",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	searchValue := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(context.Background(), searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", searchValue.SearchValue())
	}

}

func Test_Store_SearchValueUpdate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_create",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	searchValue := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(context.Background(), searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", searchValue.SearchValue())
	}

	err = store.SearchValueUpdate(context.Background(), searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", searchValue.SearchValue())
	}

	searchValue.SetSearchValue("SearchValue01Changed")
	err = store.SearchValueUpdate(context.Background(), searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "cca12edb3e3e0ba645dd430f197ff27df630fe4ec7269d31d2fab51f13d1b01a" {
		t.Fatal("Search value MUST BE 'cca12edb3e3e0ba645dd430f197ff27df630fe4ec7269d31d2fab51f13d1b01a', found: ", searchValue.SearchValue())
	}

}

func Test_Store_SearchValueFindByID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_find_by_id",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(context.Background(), value)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	valueFound, errFind := store.SearchValueFindByID(context.Background(), value.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if valueFound == nil {
		t.Fatal("Search value MUST NOT be nil")
		return
	}

	if valueFound.SourceReferenceID() != "RefId01" {
		t.Fatal("Search value reference ID MUST BE 'RefId01', found: ", valueFound.SourceReferenceID())
		return
	}

	if valueFound.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", valueFound.SearchValue())
	}

	if !strings.Contains(valueFound.SoftDeletedAt(), MAX_DATETIME) {
		t.Fatal("Search value MUST NOT be soft deleted", valueFound.SoftDeletedAt())
		return
	}
}

func Test_Store_SearchValueDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_delete",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(context.Background(), value)

	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
		return
	}

	errDelete := store.SearchValueDelete(context.Background(), value)
	if errDelete != nil {
		t.Fatal("ValueDelete Failed: " + errDelete.Error())
	}

	valueFound, errFind := store.SearchValueFindByID(context.Background(), value.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if valueFound != nil {
		t.Fatal("Search value MUST be nil")
		return
	}

	valuesFound2, errFind := store.SearchValueList(context.Background(), NewSearchValueQuery().
		SetID(value.ID()).
		SetWithSoftDeleted(true))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(valuesFound2) > 0 {
		t.Fatal("Search values MUST be 0")
		return
	}

}

func Test_Store_SearchValueSoftDelete(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_soft_delete",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(context.Background(), value)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	errDelete := store.SearchValueSoftDelete(context.Background(), value)
	if errDelete != nil {
		t.Fatal("ValueDelete Failed: " + errDelete.Error())
	}

	valueFound, errFind := store.SearchValueFindByID(context.Background(), value.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if valueFound != nil {
		t.Fatal("Search value MUST be nil")
		return
	}

	valuesFound2, errFind := store.SearchValueList(context.Background(), NewSearchValueQuery().
		SetID(value.ID()).
		SetWithSoftDeleted(true))

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(valuesFound2) > 1 {
		t.Fatal("Search values MUST NOT be 0")
		return
	}

}

func Test_Store_SearchValueSoftDeleteByID(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_soft_delete_by_id",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	if err := store.SearchValueCreate(context.Background(), value); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.SearchValueSoftDeleteByID(context.Background(), value.ID()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	values, err := store.SearchValueList(context.Background(), NewSearchValueQuery().
		SetID(value.ID()).
		SetWithSoftDeleted(true))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(values) != 1 {
		t.Fatalf("expected exactly one value with soft deleted included")
	}

	if values[0].SoftDeletedAt() == MAX_DATETIME {
		t.Fatalf("expected value to be soft deleted")
	}
}

func Test_Store_SearchEqual(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_equals",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER02",
			SearchValue: "test02@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(context.Background(), value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search(context.Background(), "test02@test.com", SEARCH_TYPE_EQUALS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 1 {
		t.Fatal("Search MUST return 1")
		return
	}

	if refsFound[0] != "USER02" {
		t.Fatal("Reference ID found MUST BE 'USER02', found: ", refsFound[0])
		return
	}
}

func Test_Store_SearchContains_NoChangeTransformer(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER021",
			SearchValue: "test021@test.com",
		},
		{
			RefID:       "USER022",
			SearchValue: "test022@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(context.Background(), value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search(context.Background(), "st02", SEARCH_TYPE_CONTAINS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 2 {
		t.Fatal("Search MUST return exactly 2 references. Returned: ", len(refsFound))
		return
	}

	if refsFound[0] != "USER021" {
		t.Fatal("Reference ID found MUST BE 'USER021', found: ", refsFound[0])
		return
	}

	if refsFound[1] != "USER022" {
		t.Fatal("Reference ID found MUST BE 'USER022', found: ", refsFound[1])
		return
	}
}

func Test_Store_Truncate(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_truncate",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	if err := store.SearchValueCreate(context.Background(), value); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.Truncate(context.Background()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	values, err := store.SearchValueList(context.Background(), NewSearchValueQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(values) != 0 {
		t.Fatalf("expected truncate to remove all rows")
	}
}

func Test_Store_SearchContains_Rot13Transformer(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &Rot13Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER021",
			SearchValue: "test021@test.com",
		},
		{
			RefID:       "USER022",
			SearchValue: "test022@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(context.Background(), value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search(context.Background(), "st02", SEARCH_TYPE_CONTAINS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 2 {
		t.Fatal("Search MUST return exactly 2 references. Returned: ", len(refsFound))
		return
	}

	if refsFound[0] != "USER021" {
		t.Fatal("Reference ID found MUST BE 'USER021', found: ", refsFound[0])
		return
	}

	if refsFound[1] != "USER022" {
		t.Fatal("Reference ID found MUST BE 'USER022', found: ", refsFound[1])
		return
	}
}

func Test_Store_SearchContains_Sha256Transformer(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER021",
			SearchValue: "test021@test.com",
		},
		{
			RefID:       "USER022",
			SearchValue: "test022@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(context.Background(), value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search(context.Background(), "st02", SEARCH_TYPE_CONTAINS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 0 {
		t.Fatal("Search MUST return exactly 0 references. Returned: ", len(refsFound))
		return
	}

	refsFound2, errFind := store.Search(context.Background(), "test022@test.com", SEARCH_TYPE_EQUALS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound2) != 1 {
		t.Fatal("Search MUST return exactly 1 references. Returned: ", len(refsFound))
		return
	}

	if refsFound2[0] != "USER022" {
		t.Fatal("Reference ID found MUST BE 'USER022', found: ", refsFound2[0])
		return
	}
}

func Test_Store_SearchContains_UniTransformer(t *testing.T) {
	db := initDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &UniTransformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER021",
			SearchValue: "test021@test.com",
		},
		{
			RefID:       "USER022",
			SearchValue: "test022@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(context.Background(), value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search(context.Background(), "st02", SEARCH_TYPE_CONTAINS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 2 {
		t.Fatal("Search MUST return exactly 2 references. Returned: ", len(refsFound))
		return
	}

	if refsFound[0] != "USER021" {
		t.Fatal("Reference ID found MUST BE 'USER021', found: ", refsFound[0])
		return
	}

	if refsFound[1] != "USER022" {
		t.Fatal("Reference ID found MUST BE 'USER022', found: ", refsFound[1])
		return
	}
}
