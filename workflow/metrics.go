package workflow

// Metrics holds workflow performance metrics.
type Metrics struct {
	ExecutionCount   ScalarMetric     `json:"execution_count"`
	SuccessCount     ScalarMetric     `json:"success_count"`
	ErrorCount       ScalarMetric     `json:"error_count"`
	AverageLatencyMs ScalarMetric     `json:"average_latency_ms"`
	LatencyOverTime  TimeSeriesMetric `json:"latency_over_time"`
	RetryRate        ScalarMetric     `json:"retry_rate"`
}

// ScalarMetric holds a single numeric metric value.
type ScalarMetric struct {
	Value float64 `json:"value"`
}

// TimeSeriesMetric holds a time series of [timestamp, value] pairs.
type TimeSeriesMetric struct {
	Value [][]float64 `json:"value"`
}

// MetricsParams holds query parameters for workflow metrics.
type MetricsParams struct {
	StartTime *string
	EndTime   *string
}
