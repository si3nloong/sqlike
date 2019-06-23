###

Functionality

TODO:

1. Logger

[ok] - [Bug] `WhereIn` and `WhereNotIn`
[ok] - Connect to `MySQL` server
[ok] - Connect to `MySQL` server with `options.Connect`
[x|gth] - Support extra options like charset and collation in connection
[ok] - Get `MySQL` server version
[ok] - Drop selected `Database`
[ok] - List all `Database`
[ok] - Support `Tag` such as `auto_increment`, `charset`, `size`, `unsigned`, `enum`, `longtext`, `generated_column`, `virtual_column`, `stored_column`
[ok] - [Feature] Truncate selected `Table`
[ok] - [Feature] List all column for selected `Table`
[ok] - [Feature] Drop selected `Table`
[ok] - [Feature] Rename `Table`
[ok] - [Feature] Check `Table` exists
[ok] - [Feature] Create single `Index` (support `unique`, `fulltext` and `spatial`)
[ok] - [Feature] Support primary key on `Migration`
[ok] - [Feature] Create multiple `Index`
[ok] - [Feature] List all `Index`
[ok] - [Feature] Drop selected `Index`
[ok] - [Feature] `Transaction` support
[ok] - [Feature] Add timeout for `Transaction`
[ok] - [Bug] Virtual column sequence in `ALTER TABLE`
[ok] - Support custom type `Key`
[ok] - Support custom type `Date`
[x] - Support custom type `GeoPoint`
[ok] - [Feature] Create `Logger`
[ok] - [Feature] Custom `JSON` encoder (w/o cover `Map` datatype)
[ok] - Custom `JSON` decoder
[ok] - [Bug] `UnmarshalJSONB` into `[]byte`
[ok] - [Feature] Single database `Migration`
[ok] - [Feature] Support generated column for `Migration` (`virtual_column` or `stored_column`)
[ok] - Support `UnsafeMigration`
[ok] - [Feature] `InsertIgnore` & `Upsert`
[x] - Set omit or setter fields on `Upsert`
[ok] - Insert single record into `Table`
[ok] - Insert multiple record into `Table`
[ok] - Retrieve single record from `Table`
[ok] - Retrieve multiple record from `Table`
[ok] - Update single record
[ok] - Modify single record by `$Key`
[ok] - Update multiple record
[ok] - Delete single record by `$Key`
[ok] - Delete multiple record
[ok] - Delete single record
[x] - Replace into (Pending for API design)
[30%] - Write testcases

[xxx] - Support `Postgres`

- Finalise Select API
- Finalise Logger API
