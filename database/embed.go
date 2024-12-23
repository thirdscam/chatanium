package db_embed

import _ "embed"

//go:embed query.sql
var DDL string
