package database

import (
	_ "embed"
)

//go:embed sqlite.schema.sql
var DbSchemaSqlite string
