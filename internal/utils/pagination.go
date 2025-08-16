package utils

// Pagination contains information about page number and limit
func Pagination(page, limit int) (skip, perPage int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50 // Default limit
	}

	skip = (page - 1) * limit
	perPage = limit
	return
}
