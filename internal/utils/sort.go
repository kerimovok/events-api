package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Sort directions
	SortAscending  = 1
	SortDescending = -1

	// Sort order values
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// BuildSortOptions converts sortBy and sortOrder into MongoDB-compatible sort options
func BuildSortOptions(sortBy, sortOrder string) bson.D {
	sortDirection := SortAscending // Default to ascending
	if sortOrder == SortOrderDesc {
		sortDirection = SortDescending
	}

	return bson.D{{Key: sortBy, Value: sortDirection}}
}
