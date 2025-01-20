package Database

import (
	"context"
	"database/sql"
	"time"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	embed "github.com/thirdscam/chatanium/database"
	"github.com/thirdscam/chatanium/src/Database/Internal"
	"github.com/thirdscam/chatanium/src/Util/Log"
)

// Database is a struct that contains the database client.
type DB struct {
	Conn    *sql.DB
	Queries *Internal.Queries
}

// Establish database connection. and, database must connected before start modules.
func (t *DB) Start() {
	ctx := context.Background()
	start := time.Now()

	// Connect to database
	db, err := sql.Open("sqlite3", ":memory:") // TODO: change to config
	if err != nil {
		Log.Error.Fatalf("Failed to connect to database: %v", err)
	}

	// Create tables
	if _, err := db.ExecContext(ctx, embed.DDL); err != nil {
		Log.Error.Fatalf("Failed to execute DDL: %v", err)
	}

	Log.Info.Printf("Connected to database. (took %s)", time.Since(start).Truncate(time.Millisecond))
	t.Queries = Internal.New(db)
	t.Conn = db
}

// Close Database connection. must be called after all modules are shutdown.
func (t *DB) Shutdown() {
	Log.Verbose.Println("Shutting down database connection...")
	if err := t.Conn.Close(); err != nil {
		Log.Error.Panicf("Cannot close database connection: %v", err)
	}
	Log.Verbose.Println("Successfully closed!")
}
