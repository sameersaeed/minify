package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP request metrics
var (
	// RequestsTotal counts total HTTP requests by method and path
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "minify_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	// RequestDuration measures request duration in seconds
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "minify_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// ResponseStatus counts responses by status code
	ResponseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "minify_http_response_status_total",
			Help: "Total number of HTTP responses by status",
		},
		[]string{"method", "path", "status"},
	)

	// URLsCreated counts the total # of Minified URLs
	URLsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "minify_urls_created_total",
			Help: "Total number of URLs minified",
		},
	)

	// URLClicks counts the total # of clicks on Minified URLs
	URLClicks = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "minify_url_clicks_total",
			Help: "Total number of URL clicks",
		},
	)

	// UsersRegistered counts the total # of registered users
	UsersRegistered = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "minify_users_registered_total",
			Help: "Total number of registered users",
		},
	)

	// ActiveUsers tracks currently active users
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "minify_active_users",
			Help: "Number of active users",
		},
	)

	// DatabaseConnections tracks active db connections
	DatabaseConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "minify_database_connections",
			Help: "Number of active database connections",
		},
	)

	// DatabaseQueryDuration measures db query duration by type
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "minify_database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)
)

// (placeholder for metric initialization, currently handled by promauto)
func Init() {}

// helpers for updating business-level metrics
func RecordURLCreated() {
	URLsCreated.Inc()
}

func RecordURLClick() {
	URLClicks.Inc()
}

func RecordUserRegistration() {
	UsersRegistered.Inc()
}

func SetActiveUsers(count float64) {
	ActiveUsers.Set(count)
}

func SetDatabaseConnections(count float64) {
	DatabaseConnections.Set(count)
}
