module github.com/maddsua/eventdb-next

go 1.23.2

replace github.com/maddsua/eventdb-next/database => ../database

require (
	github.com/99designs/gqlgen v0.17.62
	github.com/google/uuid v1.6.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/maddsua/eventdb-next/database v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/vektah/gqlparser/v2 v2.5.21
)

require (
	github.com/agnivade/levenshtein v1.2.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/sqlc-dev/sqlc v1.27.0 // indirect
)
