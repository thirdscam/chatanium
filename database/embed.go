package db_embed

import _ "embed"

//go:embed schema.sql
var DDL string
