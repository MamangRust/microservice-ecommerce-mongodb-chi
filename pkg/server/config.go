package server

import (
	"time"
)

// Config defines the common configuration for all gRPC services
type Config struct {
	ServiceName          string
	ServiceVersion       string
	Environment          string
	Port                 int
	OtelEndpoint         string
	OtelSamplingFraction float64

	// DBName adalah database MongoDB milik service ini, dipilih secara
	// eksplisit dari konstanta di pkg/database/names.go (mis. database.SalesDB).
	//
	// Wajib diisi. Server GAGAL start bila kosong (fail-hard), bukan jatuh ke
	// database default — kesalahan memilih DB harus terdeteksi saat boot, bukan
	// saat runtime menulis ke database yang salah.
	DBName string

	// MigrationPath is the directory (relative to the service workdir) that
	// holds this service's migrations. When set, migrations are applied
	// automatically at startup.
	MigrationPath string
}

// Default constants for gRPC server
const (
	DefaultMaxConcurrentConn = 1024
	DefaultWindowSize        = 16 * 1024 * 1024
	DefaultKeepaliveTime     = 20 * time.Second
	DefaultKeepaliveTimeout  = 5 * time.Second
	DefaultMinKeepaliveTime  = 5 * time.Second

	MonitoringInterval     = 30 * time.Second
	CleanupInterval        = 120 * time.Second
	CacheRefCountThreshold = 500

	ShutdownTimeout = 30 * time.Second

	RedisDialTimeout  = 5 * time.Second
	RedisReadTimeout  = 3 * time.Second
	RedisWriteTimeout = 3 * time.Second
	RedisPoolSize     = 10
	RedisMinIdleConns = 3
)
