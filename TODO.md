###

Functionality

TODO:

1. Logger
2. JSON encoder & decoder

[ok] - Connect to `MySQL` server
[ok] - Connect to `MySQL` server with `options.Connect`
[x|gth] - Support extra options like charset and collation in connection
[ok] - Get `MySQL` server version
[ok] - Drop selected `Database`
[ok] - List all `Database`
[ok] - Support `Tag` such as `auto_increment`, `size`, `unsigned`, `enum`, `longtext`, `generated_column`, `virtual_column`, `stored_column`
[ok] - Truncate selected `Table`
[ok] - List all column for selected `Table`
[ok] - Drop selected `Table`
[ok] - Check `Table` exists
[ok] - Create single `Index` (support `unique`, `fulltext` and `spatial`)
[ok] - Support primary key on `Migration`
[ok] - Create multiple `Index`
[ok] - List all `Index`
[ok] - Drop selected `Index`
[x] - Transaction support
[x] - [Bugs] Virtual column sequence in `ALTER TABLE`
[x] - Support custom type (`Key`[ok],`Date`[ok],`Point`)
[x] - Create `Logger` using `github.com/valyala/fasttemplate` // logger must be have query, arguments and milliseconds (Pending for API design)
[ok] - Custom `JSON` encoder (w/o cover `Map` datatype)
[wip] - Custom `JSON` decoder
[ok] - Single database `Migration`
[ok] - Support generated column for `Migration` (`virtual_column` or `stored_column`)
[wip] - Support `UnsafeMigration`
[ok] - `InsertIgnore` & `Upsert`
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
[20%] - Write testcases

[xxx] - Support `Postgres`

- Finalise Select API
- Finalise Logger API
