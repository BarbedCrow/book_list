package monitor

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// PgxPoolCollector exports pgxpool.Stat as Prometheus metrics.
type PgxPoolCollector struct {
	pool *pgxpool.Pool

	acquireCount      *prometheus.Desc
	acquireDuration   *prometheus.Desc
	acquiredConns     *prometheus.Desc
	idleConns         *prometheus.Desc
	totalConns        *prometheus.Desc
	maxConns          *prometheus.Desc
	emptyAcquireCount *prometheus.Desc
}

func NewPgxPoolCollector(pool *pgxpool.Pool) *PgxPoolCollector {
	return &PgxPoolCollector{
		pool: pool,
		acquireCount: prometheus.NewDesc(
			"pgxpool_acquire_count_total",
			"Total number of pool connection acquires.",
			nil, nil,
		),
		acquireDuration: prometheus.NewDesc(
			"pgxpool_acquire_duration_seconds_total",
			"Total duration of all pool acquires.",
			nil, nil,
		),
		acquiredConns: prometheus.NewDesc(
			"pgxpool_acquired_conns",
			"Number of currently acquired connections.",
			nil, nil,
		),
		idleConns: prometheus.NewDesc(
			"pgxpool_idle_conns",
			"Number of currently idle connections.",
			nil, nil,
		),
		totalConns: prometheus.NewDesc(
			"pgxpool_total_conns",
			"Total number of connections in the pool.",
			nil, nil,
		),
		maxConns: prometheus.NewDesc(
			"pgxpool_max_conns",
			"Maximum number of connections allowed.",
			nil, nil,
		),
		emptyAcquireCount: prometheus.NewDesc(
			"pgxpool_empty_acquire_count_total",
			"Total acquires that had to create a new connection.",
			nil, nil,
		),
	}
}

func (c *PgxPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.acquireCount
	ch <- c.acquireDuration
	ch <- c.acquiredConns
	ch <- c.idleConns
	ch <- c.totalConns
	ch <- c.maxConns
	ch <- c.emptyAcquireCount
}

func (c *PgxPoolCollector) Collect(ch chan<- prometheus.Metric) {
	stat := c.pool.Stat()

	ch <- prometheus.MustNewConstMetric(c.acquireCount, prometheus.CounterValue, float64(stat.AcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.acquireDuration, prometheus.CounterValue, stat.AcquireDuration().Seconds())
	ch <- prometheus.MustNewConstMetric(c.acquiredConns, prometheus.GaugeValue, float64(stat.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stat.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stat.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(stat.MaxConns()))
	ch <- prometheus.MustNewConstMetric(c.emptyAcquireCount, prometheus.CounterValue, float64(stat.EmptyAcquireCount()))
}
