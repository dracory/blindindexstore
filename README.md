# Blind Index Store

[![Tests Status](https://github.com/dracory/blindindexstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/blindindexstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/blindindexstore)](https://goreportcard.com/report/github.com/dracory/blindindexstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/blindindexstore)](https://pkg.go.dev/github.com/dracory/blindindexstore)

Blind Index Store helps you create and manage blind index database tables so you can search encrypted or tokenized data without exposing the original plaintext. It ships with reference transformers, query builders, and a high-level API for CRUD operations over blind-indexed values.

## Features

- **Turnkey store** Instantiates with `NewStore()` and optional auto-migration support.
- **Flexible transformers** Implements `TransformerInterface` so you can plug in hashing, reversible, or custom blinding strategies.
- **Rich querying** Supports equals, contains, starts-with, and ends-with search semantics (subject to transformer capabilities).
- **Soft deletes** Keeps records retrievable when you include soft-deleted values in queries.
- **Test coverage** Includes unit tests that exercise the store, query builder, and sample transformers.

## Installation

- **Go version** Requires Go 1.24+ (see `go.mod`).
- **Module**

```bash
go get github.com/dracory/blindindexstore
```

## Quick start

```go
package main

import (
    "database/sql"
    "log"

    blindindexstore "github.com/dracory/blindindexstore"
    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "file:blindindex.db?cache=shared&mode=rwc")
    if err != nil {
        log.Fatal(err)
    }

    store, err := blindindexstore.NewStore(blindindexstore.NewStoreOptions{
        DB:                 db,
        TableName:          "blindindex_emails",
        AutomigrateEnabled: true,
        Transformer:        &blindindexstore.Sha256Transformer{},
    })
    if err != nil {
        log.Fatal(err)
    }

    value := blindindexstore.NewSearchValue().
        SetSourceReferenceID("USER01").
        SetSearchValue("user01@test.com")

    if err := store.SearchValueCreate(value); err != nil {
        log.Fatal(err)
    }

    refs, err := store.Search("user01@test.com", blindindexstore.SEARCH_TYPE_EQUALS)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Refs found: %v", refs)
}
```

## Store configuration

- **`NewStoreOptions`**
  - **`TableName`** Name of the table that stores blind index entries. Required.
  - **`DB`** Any `*sql.DB` connection. Required.
  - **`DbDriverName`** Optional override. Defaults to `sb.DatabaseDriverName(DB)`.
  - **`AutomigrateEnabled`** Automatically creates the table using `sqlTableCreate()`.
  - **`DebugEnabled`** Enables SQL logging.
  - **`Transformer`** Any `TransformerInterface`. Required.

- **Automigration** Call `NewStore` with `AutomigrateEnabled: true` or invoke `store.AutoMigrate()` manually.
- **Debug logging** Toggle at runtime with `store.EnableDebug(true)`.

## Working with search values

- **Create** Use `NewSearchValue()` to build a value, then `store.SearchValueCreate()`.
- **Update** Mutate fields and call `store.SearchValueUpdate()`; the transformer re-runs when the search value changes.
- **List** Build a query with `NewSearchValueQuery()` and pass it to `store.SearchValueList()`.
- **Soft delete** Use `store.SearchValueSoftDelete()` or `SearchValueSoftDeleteByID()` and include soft-deleted rows via `SetWithSoftDeleted(true)`.
- **Hard delete** Call `store.SearchValueDelete()` or `SearchValueDeleteByID()`.

`SearchValue` instances carry metadata columns (`created_at`, `updated_at`, `soft_deleted_at`, `metas`) and helper methods for managing JSON-encoded metas.

## Searching

- **Search API** `store.Search(searchTerm, searchType)` returns matching source reference IDs.
- **Search types**
  - `SEARCH_TYPE_EQUALS`
  - `SEARCH_TYPE_CONTAINS`
  - `SEARCH_TYPE_STARTS_WITH`
  - `SEARCH_TYPE_ENDS_WITH`
- **Transformer awareness** The transformer decides what search types are meaningful. Deterministic hashes typically support equality only; reversible or passthrough transformers can support partial matches.

To build more complex queries (limits, ordering, soft delete filters), compose a `SearchValueQuery` and pass it to `store.SearchValueList()`.

## Custom transformers

- **Interface** Implement a single method: `Transform(string) string`.
- **Examples included** `NoChangeTransformer`, `Rot13Transformer`, `Sha256Transformer`, `UniTransformer` in `transformer_interface.go`.
- **Recommendations**
  - Use cryptographically strong, deterministic functions (for equality searches).
  - If supporting partial searches, document the associated trade-offs clearly.

```go
type BcryptPrefixTransformer struct{}

func (t *BcryptPrefixTransformer) Transform(v string) string {
    hashed := bcrypt.Sum([]byte(v))
    return hex.EncodeToString(hashed[:8])
}
```

Plug your transformer into `NewStoreOptions{ Transformer: &BcryptPrefixTransformer{} }`.

## Testing

- **Run tests**

```bash
go test ./...
```

The suite exercises querying behaviour, transformers, soft deletes, and truncate helpers against an in-memory SQLite database (`modernc.org/sqlite`).

## License

- **Open-source** Licensed under the GNU General Public License v3. See `LICENSE`.
- **Commercial** Contact [lesichkov.co.uk/contact](https://lesichkov.co.uk/contact) for commercial licensing options.
