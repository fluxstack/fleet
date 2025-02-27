package database

import (
	"context"
	"database/sql"
	"github.com/lynx-go/x/log"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"runtime"
)

type Client struct {
	o   Config
	DB  *bun.DB
	RDB *redis.Client
}

func NewClient(ctx context.Context, o Config) (*Client, error) {
	db, err := newDB(ctx, o.Database.Driver, o.Database.Source)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     o.Redis.Addr,
		Password: o.Redis.Password,
		DB:       o.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Client{
		DB:  db,
		RDB: rdb,
	}, nil
}

type Config struct {
	Database struct {
		Source string `json:"source"`
		Driver string `json:"driver"`
	} `json:"database"`
	Redis struct {
		Addr     string `json:"addr"`
		DB       int    `json:"db"`
		Password string `json:"password"`
	} `json:"redis"`
}

func newDB(ctx context.Context, driver string, source string) (*bun.DB, error) {
	sqldb, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	db := bun.NewDB(sqldb, mysqldialect.New())
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.InfoContext(ctx, "MySQL connection created", "source", source, "driver", driver)
	return db, nil
}
