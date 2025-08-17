package utils

const (
	// Default pagination values
	DefaultPage  = 1
	DefaultLimit = 50
)

// Pagination contains information about page number and limit
func Pagination(page, limit int) (skip, perPage int) {
	if page < DefaultPage {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}

	skip = (page - DefaultPage) * limit
	perPage = limit
	return
}
