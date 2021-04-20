# Change Log

## 2021 Apr 20

- Allow user to define namespace, override the size, charset and collate for `*types.Key`
- Introduce new api for collation 
- `NewNameKey` now using `ksuid` instead of `uuid`

## 2020 Jan 31

### Fixed :bug:

- Fix `ModifyOne` should failed when affected rows is zero
- Add `primary_key` tag support for `ModifyOne`

## 2020 Jan 29

### Feature :pill:

- Add support for `primary_key` tag

## 2019 Jul 29

- `jsonb` unmarshal bug on data type array `[]`

