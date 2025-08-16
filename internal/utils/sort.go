package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

// BuildSortOptions converts sortBy and sortOrder into MongoDB-compatible sort options
func BuildSortOptions(sortBy, sortOrder string) bson.D {
	sortDirection := 1 // Default to ascending
	if sortOrder == "desc" {
		sortDirection = -1
	}

	return bson.D{{Key: sortBy, Value: sortDirection}}
}
