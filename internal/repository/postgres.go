package repository

import (
	"context"
	"fmt"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	poolMaxConns          = 20
	poolMinConns          = 4
	poolMaxConnLifetime   = time.Minute
	poolMaxConnIdleTime   = 5 * time.Second
	poolHealthCheckPeriod = 3 * time.Second
)

func NewTransactionBalancer(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s  password=%s host=%s port=%s pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s pool_health_check_period=%s",
		cfg.Username, cfg.DbName, cfg.Password, cfg.Host, cfg.Port, poolMaxConns, poolMinConns, poolMaxConnLifetime, poolMaxConnIdleTime, poolHealthCheckPeriod)

	c, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("cant parse pool config: %s", err.Error())
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("cant configure pgxpool: %s", err.Error())
	}
	return pool, nil
}
