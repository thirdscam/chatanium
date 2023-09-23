package Database

import (
	"context"
	"os"

	"antegr.al/chatanium-bot/v1/src/Log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	Conn *pgxpool.Pool
}

func (t *Pool) Get() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		Log.Error.Fatal(err)
	}
	t.Conn = pool
}

func (t *Pool) Close() {
	t.Conn.Close()
}
