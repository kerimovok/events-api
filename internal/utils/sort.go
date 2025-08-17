package utils

import (
	"events-api/internal/constants"

	"go.mongodb.org/mongo-driver/bson"
)

// BuildSortOptions converts sortBy and sortOrder into MongoDB-compatible sort options
func BuildSortOptions(sortBy, sortOrder string) bson.D {
	sortDirection := constants.SortAscending // Default to ascending
	if sortOrder == constants.SortOrderDesc {
		sortDirection = constants.SortDescending
	}

	return bson.D{{Key: sortBy, Value: sortDirection}}
}
