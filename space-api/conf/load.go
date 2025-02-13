package conf

import _ "embed"

//go:embed sqlite.docs.sql
var Sqlite3CreateDocIndexSQLStr string
