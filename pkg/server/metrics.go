package server

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// PoolMetrics exposes MongoDB connection pool and Redis pool gauges using
// the global OTel meter provider. Values are refreshed periodically by the
// server monitoring task.
type PoolMetrics struct {
	mongoTotalConns    metric.Int64Gauge
	mongoIdleConns     metric.Int64Gauge
	mongoAcquiredConns metric.Int64Gauge
	mongoMaxConns      metric.Int64Gauge

	redisTotalConns metric.Int64Gauge
	redisIdleConns  metric.Int64Gauge
	redisHits       metric.Int64Gauge
	redisMisses     metric.Int64Gauge
}

func NewPoolMetrics(serviceName string) (*PoolMetrics, error) {
	meter := otel.Meter(serviceName)

	mongoTotalConns, err := meter.Int64Gauge(
		"mongodb_pool_total_conns",
		metric.WithDescription("Total number of connections in the MongoDB pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	mongoIdleConns, err := meter.Int64Gauge(
		"mongodb_pool_idle_conns",
		metric.WithDescription("Number of idle connections in the MongoDB pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	mongoAcquiredConns, err := meter.Int64Gauge(
		"mongodb_pool_acquired_conns",
		metric.WithDescription("Number of currently acquired connections in the MongoDB pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	mongoMaxConns, err := meter.Int64Gauge(
		"mongodb_pool_max_conns",
		metric.WithDescription("Maximum number of connections allowed in the MongoDB pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisTotalConns, err := meter.Int64Gauge(
		"redis_pool_total_conns",
		metric.WithDescription("Total number of connections in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisIdleConns, err := meter.Int64Gauge(
		"redis_pool_idle_conns",
		metric.WithDescription("Number of idle connections in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisHits, err := meter.Int64Gauge(
		"redis_pool_hits_total",
		metric.WithDescription("Total number of times a free connection was found in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisMisses, err := meter.Int64Gauge(
		"redis_pool_misses_total",
		metric.WithDescription("Total number of times a free connection was not found in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &PoolMetrics{
		mongoTotalConns:    mongoTotalConns,
		mongoIdleConns:     mongoIdleConns,
		mongoAcquiredConns: mongoAcquiredConns,
		mongoMaxConns:      mongoMaxConns,
		redisTotalConns: redisTotalConns,
		redisIdleConns:  redisIdleConns,
		redisHits:       redisHits,
		redisMisses:     redisMisses,
	}, nil
}

func (m *PoolMetrics) RecordMongoStats(ctx context.Context, total, idle, acquired, max int64) {
	m.mongoTotalConns.Record(ctx, total)
	m.mongoIdleConns.Record(ctx, idle)
	m.mongoAcquiredConns.Record(ctx, acquired)
	m.mongoMaxConns.Record(ctx, max)
}

func (m *PoolMetrics) RecordRedisStats(ctx context.Context, total, idle, hits, misses int64) {
	m.redisTotalConns.Record(ctx, total)
	m.redisIdleConns.Record(ctx, idle)
	m.redisHits.Record(ctx, hits)
	m.redisMisses.Record(ctx, misses)
}
