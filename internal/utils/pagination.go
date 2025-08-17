package utils

import "events-api/internal/constants"

// Pagination contains information about page number and limit
func Pagination(page, limit int) (skip, perPage int) {
	if page < constants.DefaultPage {
		page = constants.DefaultPage
	}
	if limit < 1 {
		limit = constants.DefaultLimit
	}

	skip = (page - constants.DefaultPage) * limit
	perPage = limit
	return
}
