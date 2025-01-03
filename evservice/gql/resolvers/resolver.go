package resolvers

import sqliteops "github.com/maddsua/eventdb-next/database/operations/sqlite"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *sqliteops.Queries
}
