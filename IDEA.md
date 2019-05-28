```go
	subQuery := sql.Table("User").
		Select("Status").
		GroupBy("Status")
	sql.Select(sql.Func("Count", sql.Col("a1"), "hats"))
	filter := sql.Table("Merchant").
		WhereIn("Status", subQuery)
	db.GetOne(filter).Decode(&model)

	db.InsertOne(
		&model,
		options.InsertOne()
		.Omit("a1", "a2", "a3")
		.SetUpsert(true)
	)

	db.InsertMany(
		&models,
		options.InsertOne()
		.Omit("a1", "a2", "a3")
		.SetUpsert(true)
	)
```

### API Should Have

- Support **JSON**
- Support **SubQuery**
- Support **Marshal** interface
- Tag should include :

1. `$Key` | `$ID` (primary key)
2. **-** (skip)
   <!-- 3. **charset:^** (UTF8 etc) -->
3. **virtual:^** (virtual column)
4. **unsigned** (only for float32 or float64)
5. **enum:FIRST|SECOND|THIRD** (enum)
6. **unique**
7. **generated**
8. **auto_increment**
