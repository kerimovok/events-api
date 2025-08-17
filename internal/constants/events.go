package constants

import "time"

// Event API constants
const (
	// Default query parameters
	DefaultSortBy    = "created_at"
	DefaultSortOrder = "asc"

	// Query parameter names
	ParamPage       = "page"
	ParamLimit      = "limit"
	ParamSortBy     = "sortBy"
	ParamSortOrder  = "sortOrder"
	ParamGroupBy    = "groupBy"
	ParamAggregates = "aggregates"
	ParamInterval   = "interval"
	ParamFilters    = "filters"

	// Default aggregation values
	DefaultAggregates = "count"
	DefaultInterval   = "day"

	// Context timeout
	QueryTimeout = 30 * time.Second

	// Pagination constants
	DefaultPage  = 1
	DefaultLimit = 50

	// Sort constants
	SortAscending  = 1
	SortDescending = -1
	SortOrderAsc   = "asc"
	SortOrderDesc  = "desc"

	// Collection names
	EventsCollection = "events"

	// Aggregation operations
	AggregationCount = "count"
	AggregationSum   = "sum"
	AggregationAvg   = "avg"

	// Time intervals
	IntervalHour  = "hour"
	IntervalDay   = "day"
	IntervalWeek  = "week"
	IntervalMonth = "month"

	// Time interval formats (for MongoDB dateToString)
	TimeFormatHour  = "%Y-%m-%d-%H"
	TimeFormatDay   = "%Y-%m-%d"
	TimeFormatWeek  = "%Y-%U"
	TimeFormatMonth = "%Y-%m"
)
